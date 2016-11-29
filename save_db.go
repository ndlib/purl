package main

import "database/sql"

type dbObj struct {
	db *sql.DB // store pointer to sql database
}

func NewDBFileStore(db *sql.DB) *purldb {
	return &purldb{db: db}
}
