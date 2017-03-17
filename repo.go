package main

import (
	"net/http"
)

// A Repository defines the actions we need to do against the
// PURL database. It is an interface so we can use either a MySQL or
// Postgres backend (or an in-memory one for testing).
type Repository interface {
	AllPurls() []Purl

	AllRepos() []RepoObj

	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) Purl

	FindRepoObj(id int) RepoObj

	FindQuery(query string) []RepoObj

	CreatePurl(t Purl)

	CreateRepo(t RepoObj)

	LogRecordAccess(vars *http.Request, repoID int, purlID int)
}
