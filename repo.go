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

// A Repository stores the purls we know about.
type Repository interface {
	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) Purl

	//
	FindQuery(query string) []RepoObj

	AllPurls() []Purl

	CreatePurl(t Purl)
}

type purldb struct {
	db *sql.DB // store pointer to sql database
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
	return
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

// type Purl struct {
// 	Id            int       `json:"id"`
// 	repo_obj_id   string    `json:"repo_obj_id"`
// 	access_count  int       `json:"access_count"`
// 	last_accessed time.Time `json:"last_accessed"`
// 	source_app    string    `json:"source_app"`
// 	date_created  time.Time `json:"date_created"`
// }
// refernce to purldb
func (sq *purldb) AllPurls() []Purl {
	var result []Purl
	qstring := "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl"
	rows, err := sq.db.Query(qstring)
	if err != nil {
		log.Printf("Error getting all purls: %s", err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var temp_purl Purl
		var last_accessed mysql.NullTime
		var source_app sql.NullString
		err := rows.Scan(&temp_purl.Id, &temp_purl.Repo_obj_id, &temp_purl.Access_count, &last_accessed, &source_app, &temp_purl.Date_created)
		if err != nil {
			log.Printf("Scan not succeded: %s", err)
		}
		fmt.Println(len(result))
		if last_accessed.Valid {
			temp_purl.Last_accessed = last_accessed.Time
		}
		if source_app.Valid {
			temp_purl.Source_app = source_app.String
		}
		result = append(result, temp_purl)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error on rows scan: %s", err)
		return result
	}
	return result
}

func (sq *purldb) FindPurl(id int) Purl {
	result := Purl{}
	qstring := "select purl_id, repo_object_id, access_count, last_accessed, source_app, date_created from purl where purl_id = ?"
	row, err := sq.db.Query(qstring, id)
	if err != nil {
		log.Printf("Error getting all purls: %s", err.Error())
		return result
	}
	defer row.Close()
	var purl_id int
	var repo_obj_id string
	var access_count int
	var last_accessed time.Time
	var source_app string
	var date_created time.Time
	if err := row.Scan(&purl_id, &repo_obj_id, &access_count, &last_accessed, &source_app, &date_created); err != nil {
		log.Printf("Error scanning rows: %s", err.Error())
		return result
	}
	result = Purl{Id: purl_id, Repo_obj_id: repo_obj_id, Access_count: access_count, Last_accessed: last_accessed, Source_app: source_app, Date_created: date_created}
	// return empty if not found
	return result
}

func (sq *purldb) FindQuery(query string) []RepoObj {
	var result []RepoObj
	qstring := "select filename, url, date_added, add_source_ip, date_modified, information from repo_object where repo_object.information = ?"
	rows, err := sq.db.Query(qstring, query)
	if err != nil {
		log.Printf("Error getting all purls: %s", err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var filename string
		var url string
		var date_added time.Time
		var add_source_ip string
		var date_modified time.Time
		var information string
		if err := rows.Scan(&filename, &url, &date_added, &add_source_ip, &date_modified, &information); err != nil {
			log.Fatal(err)
			return result
		}
		result = append(result, RepoObj{Filename: filename, Url: url, Date_added: date_added, Add_source_ip: add_source_ip, Date_modified: date_modified, Information: information})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return result
	}
	return result
}

func (sq *purldb) CreatePurl(t Purl) {
	qstring := "INSERT INTO purl (repo_object_id, access_count, last_accessed, date_created) VALUES (?,?,?,?)"
	_, err := sq.db.Query(qstring, t.Repo_obj_id, t.Access_count, t.Last_accessed, t.Date_created)
	if err != nil {
		log.Printf("Could not insert into database: %s", err.Error())
		return
	}
	return
}

// func (mr *purldb) DestroyPurl(id int) error {
// 	for i, t := range mr.purls {
// 		if t.Id == id {
// 			mr.purls = append(mr.purls[:i], mr.purls[i+1:]...)
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("Could not find Purl with id of %d to delete", id)
// }
