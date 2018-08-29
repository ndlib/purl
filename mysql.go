package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

// store pointer to sql database
type mysqlDB struct {
	db *sql.DB
}

// NewMySQL returns a Repository backed by a MySQL database, as determined
// by the connection string. An error is returned if any problems are run into.
func NewMySQL(conn string) (Repository, error) {
	conn += "?parseTime=true"
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	return &mysqlDB{db: db}, nil
}

// AllPurls returns a list of every purl in the database.
func (sq *mysqlDB) AllPurls() []Purl {
	const query = `SELECT
		p.purl_id, p.repo_object_id, p.access_count, p.last_accessed, p.date_created, r.filename, r.url, r.information
		FROM purl as p left join repo_object as r on p.repo_object_id = r.repo_object_id;`

	var result []Purl
	rows, err := sq.db.Query(query)
	if err != nil {
		log.Println("Error getting all purls:", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		result = append(result, scanPurl(rows))
	}
	if err := rows.Err(); err != nil {
		log.Println("Error on rows scan:", err)
		return result
	}
	return result
}

// scanPurl reads the current row and returns a Purl.
func scanPurl(rows *sql.Rows) Purl {
	var p Purl
	var lastAccessed mysql.NullTime
	var information sql.NullString
	err := rows.Scan(&p.ID, &p.RepoID, &p.AccessCount, &lastAccessed, &p.DateCreated, &p.Filename, &p.URL, &information)
	if err != nil {
		log.Println("Scan:", err)
		return p
	}
	if lastAccessed.Valid {
		p.LastAccessed = lastAccessed.Time
	}
	if information.Valid {
		p.Information = information.String
	}
	return p
}

// FindPurl returns a single Purl record as given by id. The bool is true if the asked-for
// record was found, false if it was not found.
func (sq *mysqlDB) FindPurl(id int) (Purl, bool) {
	rows, err := sq.db.Query(`
		SELECT p.purl_id, p.repo_object_id, p.access_count, p.last_accessed, p.date_created, r.filename, r.url, r.information
		FROM purl as p left join repo_object as r on p.repo_object_id = r.repo_object_id
		WHERE purl_id = ?
		LIMIT 1;`,
		id)
	if err != nil {
		log.Println("Purl", id, ":", err)
		return Purl{}, false
	}
	defer rows.Close()
	for rows.Next() {
		return scanPurl(rows), true
	}
	return Purl{}, false
}

// FIND FOR REPO INFORMATION
func (sq *mysqlDB) FindQuery(query string) []Purl {
	var result []Purl
	if query == "" {
		return sq.AllPurls()
	}
	const qstring = `SELECT p.purl_id, p.repo_object_id, p.access_count, p.last_accessed, p.date_created, r.filename, r.url, r.information
		FROM purl as p left join repo_object as r on p.repo_object_id = r.repo_object_id
		WHERE r.information LIKE ?;`
	rows, err := sq.db.Query(qstring, "%"+query+"%")
	if err != nil {
		log.Println("FindQuery:", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		result = append(result, scanPurl(rows))
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return result
}

func (sq *mysqlDB) SummaryStats() Stats {
	return Stats{}
}

// LOGS ACCESS TO THE DATABASE
func (sq *mysqlDB) LogAccess(r *http.Request, purl Purl) {
	const upstring = `INSERT INTO object_access
	(date_accessed, ip_address, host_name, referer, user_agent, request_method, path_info, repo_object_id, purl_id)
	VALUES
	(now(),?,?,?,?,?,?,?,?)`

	_, err := sq.db.Exec(
		upstring,
		r.RemoteAddr,
		r.Host,
		r.Referer(),
		r.UserAgent(),
		r.Method,
		r.URL.Path,
		purl.RepoID,
		purl.ID,
	)
	if err != nil {
		log.Println("LogAccess:", err)
	}

	const purlUpdate = `UPDATE purl
		SET access_count=access_count+1, last_accessed=now()
		WHERE purl_id = ?
		LIMIT 1`
	_, err = sq.db.Exec(purlUpdate, purl.ID)
	if err != nil {
		log.Println("LogAccess:", err)
	}
}
