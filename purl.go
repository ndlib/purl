package main

import "time"

type Purl struct {
	Id            int       `json:"id"`
	Repo_obj_id   string    `json:"repo_obj_id"`
	Access_count  int       `json:"access_count"`
	Last_accessed time.Time `json:"last_accessed"`
	Source_app    string    `json:"source_app"`
	Date_created  time.Time `json:"date_created"`
}

type Purls []Purl

type RepoObj struct {
	Id            int       `json:"id"`
	Filename      string    `json:"filename"`
	Url           string    `json:"filename"`
	Date_added    time.Time `json:"date_added"`
	Add_source_ip string    `json:"add_source_ip"`
	Date_modified time.Time `json:"date_modified"`
	Information   string    `json:"information"`
}

type Repos []RepoObj
