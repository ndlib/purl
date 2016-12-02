package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func QueryDB(sq *purldb, id int) (*sql.Rows, error) {
	var qstring string
	if id == -1 {
		qstring = "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl"
		return sq.db.Query(qstring)
	} else {
		qstring = "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl where purl_id = ?"
		return sq.db.Query(qstring, id)
	}
}

// type Purl struct {
// 	Id            int       `json:"id"`
// 	repo_obj_id   string    `json:"repo_obj_id"`
// 	access_count  int       `json:"access_count"`
// 	last_accessed time.Time `json:"last_accessed"`
// 	source_app    string    `json:"source_app"`
// 	date_created  time.Time `json:"date_created"`
// }
// refernce to purldb

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

func CreatePurlDB(sq *purldb, purl Purl) sql.Result {
  var qstring string
  qstring = "INSERT INTO purl (purl_id, repo_object_id, last_accessed, source_app, date_created) VALUES ($1, $2, $3, $4, $5, $6)"
  result, err := sq.db.Exec(
    qstring,
    purl.Id, purl.Repo_obj_id, purl.Last_accessed, purl.Source_app, purl.Date_created
  )
  if err != nil {
    log.Printf("Error creating purl: %s", err.Error())
    return result
  }
  return result
}

// type RepoObj struct {
// 	Id            int       `json:"id"`
// 	Filename      string    `json:"filename"`
// 	Url           string    `json:"filename"`
// 	Date_added    time.Time `json:"date_added"`
// 	Add_source_ip string    `json:"add_source_ip"`
// 	Date_modified time.Time `json:"date_modified"`
// 	Information   string    `json:"information"`
// }

func ScanRepoDB(rows *sql.Rows) RepoObj {
	var temp_repo RepoObj
	var date_modified mysql.NullTime
	var information sql.NullString
	err := rows.Scan(&temp_repo.Id, &temp_repo.Filename, &temp_repo.Url, &temp_repo.Date_added, &temp_repo.Add_source_ip, &date_modified, &information)
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

func CreateRepoDB(sq *purldb, repo Repo) sql.Result {
  var qstring string
  qstring = "INSERT INTO repo_object (repo_object_id, filename, url, date_added, add_source_ip, date_modified, information) VALUES ($1, $2, $3, $4, $5, $6, $7)"
  result, err := sq.db.Exec(
    qstring,
    repo.Id, repo.Filename, repo.Url, repo.Date_added, repo.Add_source_ip, repo.Date_modified, repo.Information
  )
  if err != nil {
    log.Printf("Error creating repo: %s", err.Error())
    return result
  }
  return result
}
