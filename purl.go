package main

import "time"

type Purl struct {
	Id            int       `json:"id"`
	Repo_obj_id   string    `json:"repo_obj_id"`
	Access_count  int       `json:"access_count"`
	Last_accessed time.Time `json:"last_accessed"`
	Repo_url    string      `json:"repo_url"`
	Purl_url    string      `json:"purl_url"`
	Date_created  time.Time `json:"date_created"`
	File_name 		string		`json:"file_name"`
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
