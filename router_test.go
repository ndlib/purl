// +build router

package main

import (
	"testing"
	"sync"
	"httptest"

)

var (
	source *purldb
)

// A MemoryRepo is a Repository that keeps everything in memory.
// It is mostly useful for testing.
type memoryRepo struct {
	m sync.RWMutex // protects everything below

	// last ID minted
	currentID int

	// list of Purl objects
	purls []Purl

	// list of repository resources
	repos []Repo
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

func (mr *memoryRepo) FindRepos(id int) RepoObj {
	mr.m.RLock()
	defer mr.m.RUnlock()
	var ret RepoObj
	for _, r := range mr.repos {
		if r.Id = id {
			return r
		}
	}
	return RepoObj{}
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

func TestRouterProxy(t *testing.T){
	handler := &PurlShowFile{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.Url)
}

func init() {

}




















