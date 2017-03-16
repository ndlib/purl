package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	datasource Repository
)

func main() {
	log.Println("===== Starting Repopurl")
	dbLocation := os.Getenv("MYSQL_CONNECTION")
	if dbLocation == "" {
		dbLocation = "/repopurl"
		log.Println("Using default database", dbLocation)
	}
	// don't log dblocation otherwise, since it probably contains passwords
	datasource = NewDBSource(dbLocation)

	router := NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
