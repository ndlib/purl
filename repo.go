package main

import "fmt"

var currentId int

var purls Purls

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
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
