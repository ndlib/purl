package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type dbObj struct {
	db *sql.DB // store pointer to sql database
}

func NewDBFileStore(db *sql.DB) *purldb {
	return &purldb{db: db}
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
		return result
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
