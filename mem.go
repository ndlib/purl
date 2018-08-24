package main

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// memoryRepo implements a Repository that is kept entirely in memory.
// It is expected to be useful mostly for testing.
type memoryRepo struct {
	m         sync.RWMutex // protects everything below
	currentID int          // last ID minted
	purls     []Purl       // list of Purl objects
	history   []Purl       // access history
}

func (mr *memoryRepo) AllPurls() []Purl {
	mr.m.RLock()
	defer mr.m.RUnlock()
	return mr.purls[:]
}

func (mr *memoryRepo) FindPurl(id int) (Purl, bool) {
	mr.m.RLock()
	defer mr.m.RUnlock()
	for _, p := range mr.purls {
		if p.ID == id {
			return p, true
		}
	}
	// return empty if not found
	return Purl{}, false
}

func (mr *memoryRepo) FindQuery(query string) []Purl {
	mr.m.RLock()
	defer mr.m.RUnlock()
	var result []Purl
	for _, p := range mr.purls {
		if strings.Contains(p.Information, query) {
			result = append(result, p)
		}
	}
	return result
}

func (mr *memoryRepo) SummaryStats() Stats {
	return Stats{}
}

func (mr *memoryRepo) CreatePurl(p Purl) {
	mr.m.Lock()
	defer mr.m.Unlock()
	if p.ID == 0 {
		mr.currentID++
		p.ID = mr.currentID
	}
	mr.purls = append(mr.purls, p)
}

func (mr *memoryRepo) LogAccess(vars *http.Request, purl Purl) {
	mr.m.Lock()
	defer mr.m.Unlock()
	mr.history = append(mr.history, purl)
	for i := range mr.purls {
		if mr.purls[i].ID == purl.ID {
			mr.purls[i].LastAccessed = time.Now()
			mr.purls[i].AccessCount++
			return
		}
	}
}
