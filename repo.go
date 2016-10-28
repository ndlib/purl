package main

import (
	"fmt"
	"strings"
)

// A Repository stores the purls we know about.
type Repository interface {
	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) Purl

	//
	FindQuery(query string) []RepoObj

	AllPurls() []Purl

	CreatePurl(t Purl) Purl
}

type memoryRepo struct {
	currentID int
	purls     Purls
	repos     Repos
}

func (mr *memoryRepo) AllPurls() []Purl {
	return mr.purls
}

func (mr *memoryRepo) FindPurl(id int) Purl {
	for _, t := range mr.purls {
		if t.Id == id {
			return t
		}
	}
	// return empty if not found
	return Purl{}
}

func (mr *memoryRepo) FindQuery(query string) []RepoObj {
	var ret []RepoObj
	for _, q := range mr.repos {
		if strings.Contains(q.information, query) {
			ret = append(ret, q)
		}
	}
	return ret
}

//this is bad, I don't think it passes race condtions
func (mr *memoryRepo) CreatePurl(t Purl) Purl {
	mr.currentID += 1
	t.Id = mr.currentID
	mr.purls = append(mr.purls, t)
	return t
}

func (mr *memoryRepo) DestroyPurl(id int) error {
	for i, t := range mr.purls {
		if t.Id == id {
			mr.purls = append(mr.purls[:i], mr.purls[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Purl with id of %d to delete", id)
}
