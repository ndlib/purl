package main

import (
	"time"
)

// A Purl represents a single redirect in the database.
type Purl struct {
	ID           int       `json:"id"`
	RepoObjID    string    `json:"repo_obj_id"`
	AccessCount  int       `json:"access_count"`
	LastAccessed time.Time `json:"last_accessed"`
	SourceApp    string    `json:"source_app"`
	DateCreated  time.Time `json:"date_created"`
}

// A RepoObj is a braindead structure needed because we are mirroring
// how the PURLs are stored in the database. For the most part there
// is always a one-to-one relationship between a Purl and a RepoObj.
type RepoObj struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	URL          string    `json:"URL"`
	DateAdded    time.Time `json:"date_added"`
	AddSourceIP  string    `json:"add_source_ip"`
	DateModified time.Time `json:"date_modified"`
	Information  string    `json:"information"`
}
