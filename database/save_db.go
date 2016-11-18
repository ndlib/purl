package main

import (
  "datebase/sql"
  "log"
  "time"
)

type dbObj struct {
	db *sql.DB // store pointer to sql database
}

func NewDBFileStore(db *sql.DB) purl_db {
  return &dbObj(DB: db)
}
