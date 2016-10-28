package main

import (
	"log"
	"net/http"
)

var (
	datasource Repository
)

func main() {

	repo := &memoryRepo{}
	// Give us some seed data
	repo.CreatePurl(Purl{source_app: "Library"})
	repo.CreatePurl(Purl{source_app: "Code School"})

	datasource = repo

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
