package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
)

// store pointer to sql database
type purldb struct {
	db *sql.DB
}

// A MemoryRepo is a Repository that keeps everything in memory.
// It is mostly useful for testing.
type memoryRepo struct {
	m sync.RWMutex // protects everything below

	// last ID minted
	currentID int

	// list of Purl objects
	purls Purls

	// list of repository resources
	repos Repos
}

// A Repository stores the purls we know about.
type Repository interface {
	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) Purl

	FindQuery(query string) []RepoObj

	AllPurls() []Purl

	CreatePurl(t Purl)
}

func updateWait(wait int) int {
	wait *= 2
	if wait > 300 {
		wait = 300
	}
	return wait
}

func NewDBSource(db *sql.DB, mysqlconn string) (*sql.DB, *purldb) {
	var err error
	db, err = sql.Open("mysql", mysqlconn)
	if err != nil {
		var wait = 1
		for err != nil {
			log.Printf("Error opening database: %s", err.Error())
			db, err = sql.Open("mysql", mysqlconn)
			time.Sleep(time.Duration(wait) * time.Second)
			wait = updateWait(wait)
		}
	}
	db.Ping()
	return db, &purldb{db: db}
}

func (sq *purldb) queryDB(id int) (*sql.Rows, error) {
	var qstring string
	if id == -1 {
		qstring = "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl"
		return sq.db.Query(qstring)
	} else {
		qstring = "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl where purl_id = ?"
		return sq.db.Query(qstring, id)
	}
}

func ScanPurlDB(rows *sql.Rows) Purl {
	var temp_purl Purl
	var last_accessed mysql.NullTime
	var source_app sql.NullString
	err := rows.Scan(&temp_purl.Id, &temp_purl.Repo_obj_id, &temp_purl.Access_count, &last_accessed, &source_app, &temp_purl.Date_created)
	if err != nil {
		log.Printf("Scan not succeded: %s", err)
		return temp_purl
	}
	if last_accessed.Valid {
		temp_purl.Last_accessed = last_accessed.Time
	}
	if source_app.Valid {
		temp_purl.Source_app = source_app.String
	}
	return temp_purl
}

func (sq *purldb) createPurlDB(purl Purl) sql.Result {
	var qstring string
	qstring = `INSERT INTO purl
	(purl_id, repo_object_id, last_accessed, source_app, date_created)
	VALUES
	(?, ?, ?, ?, ?)`
	result, err := sq.db.Exec(
		qstring,
		purl.Id, purl.Repo_obj_id, purl.Last_accessed, purl.Source_app, purl.Date_created,
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
		return result
	}
	return result
}

func ScanRepoDB(rows *sql.Rows) RepoObj {
	var temp_repo RepoObj
	var date_modified mysql.NullTime
	var information sql.NullString
	err := rows.Scan(&temp_repo.Id, &temp_repo.Filename, &temp_repo.Url, &temp_repo.Date_added,
		&temp_repo.Add_source_ip, &date_modified, &information)
	if err != nil {
		log.Printf("Scan not succeded: %s", err)
		return temp_repo
	}
	if date_modified.Valid {
		temp_repo.Date_modified = date_modified.Time
	}
	if information.Valid {
		temp_repo.Information = information.String
	}
	return temp_repo
}

func (sq *purldb) createRepoDB(repo RepoObj) sql.Result {
	var qstring string
	qstring = `INSERT INTO repo_object
	(repo_object_id, filename, url, date_added, add_source_ip, date_modified, information)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)`
	result, err := sq.db.Exec(
		qstring,
		repo.Id, repo.Filename, repo.Url, repo.Date_added, repo.Add_source_ip, repo.Date_modified, repo.Information,
	)
	if err != nil {
		log.Printf("Error creating repo: %s", err.Error())
		return result
	}
	return result
}

func (mr *memoryRepo) AllPurls() []Purl {
	mr.m.RLock()
	defer mr.m.RUnlock()
	return mr.purls[:]
}

func (mr *memoryRepo) FindPurl(id int) Purl {
	mr.m.RLock()
	defer mr.m.RUnlock()
	for _, t := range mr.purls {
		if t.Id == id {
			return t
		}
	}
	// return empty if not found
	return Purl{}
}

func (mr *memoryRepo) FindQuery(query string) []RepoObj {
	mr.m.RLock()
	defer mr.m.RUnlock()
	var ret []RepoObj
	for _, q := range mr.repos {
		if strings.Contains(q.Information, query) {
			ret = append(ret, q)
		}
	}
	return ret
}

func (mr *memoryRepo) CreatePurl(t Purl) {
	mr.m.Lock()
	defer mr.m.Unlock()
	mr.currentID += 1
	t.Id = mr.currentID
	mr.purls = append(mr.purls, t)
}

func (mr *memoryRepo) DestroyPurl(id int) error {
	mr.m.Lock()
	defer mr.m.Unlock()
	for i, t := range mr.purls {
		if t.Id == id {
			mr.purls = append(mr.purls[:i], mr.purls[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Purl with id of %d to delete", id)
}

// Queries the database for all purls and
// returns an empty object if there is an
// error
func (sq *purldb) AllPurls() []Purl {
	var result []Purl
	rows, err := sq.queryDB(-1)
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

func (sq *purldb) FindPurl(id int) Purl {
	result := Purl{}
	row, err := sq.queryDB(id)
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

func (sq *purldb) CreatePurl(t Purl) {
	sq.createPurlDB(t)
	return
}
