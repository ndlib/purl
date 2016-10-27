package main

import "time"

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
