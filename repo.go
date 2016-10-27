package main

import (
	"fmt"
	"strings"
)

var currentId int

var purls Purls
var repos Repos

// Give us some seed data
func init() {
	RepoCreatePurl(Purl{source_app: "Library"})
	RepoCreatePurl(Purl{source_app: "Code School"})
}

func RepoFindPurl(id int) Purl {
	for _, t := range purls {
		if t.Id == id {
			return t
		}
	}
	// return empty if not found
	return Purl{}
}

func RepoFindQuery(query string) []RepoObj {
	var ret []RepoObj
	for _, q := range repos {
		if strings.Contains(q.information, query) {
			ret = append(ret, q)
		}
	}
	return ret
}

func RepoFindPurlFile(id int, file string) Purl {
	for _, t := range purls {
		if t.Id == id {
			return t
		}
	}
	// return empty if not found
	return Purl{}
}

//this is bad, I don't think it passes race condtions
func RepoCreatePurl(t Purl) Purl {
	currentId += 1
	t.Id = currentId
	purls = append(purls, t)
	return t
}

func RepoDestroyPurl(id int) error {
	for i, t := range purls {
		if t.Id == id {
			purls = append(purls[:i], purls[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Purl with id of %d to delete", id)
}
