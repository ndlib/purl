package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
)

// store pointer to sql database
type purldb struct {
	db *sql.DB
}

func updateWait(wait int) int {
	wait *= 2
	if wait > 300 {
		wait = 300
	}
	return wait
}

// NewDBSource returns a Repository backed by a MySQL database, as determined
// by the connection string.
func NewDBSource(mysqlconn string) *purldb {
	mysqlconn = mysqlconn + "?parseTime=true"
	db, err := sql.Open("mysql", mysqlconn)
	if err != nil {
		log.Printf("Error setting up database connection: %s", err.Error())
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		var wait = 1
		for err != nil {
			log.Printf("Error opening database: %s", err.Error())
			err = db.Ping()
			time.Sleep(time.Duration(wait) * time.Second)
			wait = updateWait(wait)
		}
	}
	return &purldb{db: db}
}

func (sq *purldb) createPurlDB(purl Purl) sql.Result {
	var qstring string
	qstring = `INSERT INTO purl
	(purl_id, repo_object_id, last_accessed, source_app, date_created)
	VALUES
	(?, ?, ?, ?, ?)`
	result, err := sq.db.Exec(
		qstring,
		purl.ID, purl.RepoObjID, purl.LastAccessed, purl.SourceApp, purl.DateCreated,
	)
	if err != nil {
		log.Printf("Error creating purl: %s", err.Error())
	}
	return result
}

func (sq *purldb) destroyPurlDB(id int) sql.Result {
	var qstring string
	qstring = "delete from purl where purl_id = ?"
	result, err := sq.db.Exec(
		qstring,
		id,
	)
	if err != nil {
		log.Printf("Error creating purl: %s", err.Error())
	}
	return result
}

func (sq *purldb) createRepoDB(repo RepoObj) sql.Result {
	var qstring string
	qstring = `INSERT INTO repo_object
	(repo_object_id, filename, url, date_added, add_source_ip, date_modified, information)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)`
	result, err := sq.db.Exec(
		qstring,
		repo.ID, repo.Filename, repo.URL, repo.DateAdded, repo.AddSourceIP, repo.DateModified, repo.Information,
	)
	if err != nil {
		log.Printf("Error creating repo: %s", err.Error())
	}
	return result
}

// Queries the database for all purls and
// returns an empty object if there is an
// error

func (sq *purldb) queryDB(id int, table string, tableID string) (*sql.Rows, error) {
	if id == -1 {
		qstring := "select * from " + table
		return sq.db.Query(qstring)
	}
	qstring := "select * from " + table + " where " + tableID + " = ?"
	upstring := "UPDATE purl SET access_count = access_count + 1, last_accessed = NOW()"
	_, _ = sq.db.Query(upstring, id)
	return sq.db.Query(qstring, id)
}

// PURL OBJECT RETRIEVAL
func (sq *purldb) queryPurlDB(id int) (*sql.Rows, error) {
	return sq.queryDB(id, "purl", "purl_id")
}

// REPO OBJECT RETRIEVAL
func (sq *purldb) queryRepoDB(id int) (*sql.Rows, error) {
	return sq.queryDB(id, "repo_object", "repo_object_id")
}

// FULL PURL LISTING RETRIEVAL
func (sq *purldb) AllPurls() []Purl {
	var result []Purl
	rows, err := sq.queryPurlDB(-1)
	if err != nil {
		log.Printf("Error getting all purls: %s", err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		result = append(result, ScanPurlDB(rows))
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error on rows scan: %s", err)
		return result
	}
	return result
}

// FULL REPO RESOURCE RETRIEVAL
func (sq *purldb) AllRepos() []RepoObj {
	var result []RepoObj
	rows, err := sq.queryRepoDB(-1)
	if err != nil {
		log.Printf("Error getting all repos: %s", err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		result = append(result, ScanRepoDB(rows))
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error on rows scan: %s", err)
		return result
	}
	return result
}

// ScanPurlDB reads the current row in rows and returns a Purl.
func ScanPurlDB(rows *sql.Rows) Purl {
	var tempPurl Purl
	var lastAccessed mysql.NullTime
	var sourceApp sql.NullString
	err := rows.Scan(
		&tempPurl.ID, &tempPurl.RepoObjID, &tempPurl.AccessCount,
		&lastAccessed, &sourceApp, &tempPurl.DateCreated,
	)
	if err != nil {
		log.Printf("Scan not succeeded: %s", err)
		return tempPurl
	}
	if lastAccessed.Valid {
		tempPurl.LastAccessed = lastAccessed.Time
	}
	if sourceApp.Valid {
		tempPurl.SourceApp = sourceApp.String
	}
	return tempPurl
}

// PURL OBJECT SEARCH AND RETRIEVAL
func (sq *purldb) FindPurl(id int) Purl {
	result := Purl{}
	row, err := sq.queryPurlDB(id)
	if err != nil {
		log.Printf("Error getting purl %d purls: %s", id, err.Error())
		return result
	}
	defer row.Close()
	for row.Next() {
		result = ScanPurlDB(row)
	}
	return result
}

// ScanRepoDB reads the current row in rows and returns a corresponding RepoObj.
func ScanRepoDB(rows *sql.Rows) RepoObj {
	var tempRepo RepoObj
	var dateModified mysql.NullTime
	var information sql.NullString
	err := rows.Scan(&tempRepo.ID, &tempRepo.Filename, &tempRepo.URL, &tempRepo.DateAdded,
		&tempRepo.AddSourceIP, &dateModified, &information)
	if err != nil {
		log.Printf("Scan not succeeded: %s", err)
		return tempRepo
	}
	if dateModified.Valid {
		tempRepo.DateModified = dateModified.Time
	}
	if information.Valid {
		tempRepo.Information = information.String
	}
	return tempRepo
}

// REPO OBJECT SEARCH AND RETRIEVAL
func (sq *purldb) FindRepoObj(id int) RepoObj {
	result := RepoObj{}
	row, err := sq.queryRepoDB(id)
	if err != nil {
		log.Printf("Error getting repo %d repos: %s", id, err.Error())
		return result
	}
	defer row.Close()
	for row.Next() {
		result = ScanRepoDB(row)
	}
	return result
}

// FIND FOR REPO INFORMATION
func (sq *purldb) FindQuery(query string) []RepoObj {
	var result []RepoObj
	qstring := `SELECT filename, url, date_added, add_source_ip, date_modified, information from repo_object
							where
							repo_object.information CONTAINS ?`
	rows, err := sq.db.Query(qstring, query)
	if err != nil {
		log.Printf("Error getting all purls: %s", err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		result = append(result, ScanRepoDB(rows))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return result
	}
	return result
}

// PURL TESTING AND CREATION
func (sq *purldb) CreatePurl(t Purl) {
	sq.createPurlDB(t)
	return
}

// REPO TESTING AND CREATION
func (sq *purldb) CreateRepo(t RepoObj) {
	sq.createRepoDB(t)
	return
}

// LOGS ACCESS TO THE DATABASE
func (sq *purldb) LogRecordAccess(r *http.Request, repoID int, purlID int) {
	upstring := `INSERT INTO object_access
	(date_accessed, ip_address, host_name, referer, user_agent, request_method, path_info, repo_object_id, purl_id)
	VALUES
	(now(),?,?,?,?,?,?,?,?)`
	_, err := sq.db.Exec(
		upstring, r.RemoteAddr, r.Host, r.Referer(), r.UserAgent(),
		r.Method, r.URL.Path, repoID, purlID,
	)
	if err != nil {
		log.Printf("Problem updating access to database: %s", err.Error())
	}
}
