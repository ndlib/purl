package main

import (
	"fmt"
	"strings"
	"sync"
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
		if strings.Contains(q.information, query) {
			ret = append(ret, q)
		}
	}
	return ret
}

func (mr *memoryRepo) CreatePurl(t Purl) Purl {
	mr.m.Lock()
	defer mr.m.Unlock()
	mr.currentID += 1
	t.Id = mr.currentID
	mr.purls = append(mr.purls, t)
	return t
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
