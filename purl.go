package main

import "time"

type Purl struct {
	Id            int       `json:"id"`
	repo_obj_id   string    `json:"repo_obj_id"`
	access_count  int       `json:"access_count"`
	last_accessed time.Time `json:"last_accessed"`
	source_app    string    `json:"source_app"`
	date_created  time.Time `json:"date_created"`
}

type Purls []Purl

type RepoObj struct {
	Id            int       `json:"id"`
	filename      string    `json:"filename"`
	url           string    `json:"filename"`
	date_added    time.Time `json:"date_added"`
	add_source_ip string    `json:"add_source_ip"`
	date_modified time.Time `json:"date_modified"`
	information   string    `json:"information"`
}

type Repos []RepoObj
