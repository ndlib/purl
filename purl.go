package main

import "time"

type Purl struct {
	Id                    int         `json:"id"`
	repo_obj_id           string      `json:"repo_obj_id"`
	access_count          int         `json:"access_count"`
  last_accessed         time.Time    `json:"last_accessed"`
  source_app            string       `json:"source_app"`
	date_created          time.Time    `json:"date_created"`
}

type Purls []Purl
