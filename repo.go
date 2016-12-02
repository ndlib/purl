package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
)

// A Repository stores the purls we know about.
type Repository interface {
	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) Purl

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
	qstring := "select filename, url, date_added, add_source_ip, date_modified, information from repo_object where repo_object.information = ?"
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
