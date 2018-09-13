package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("===== Starting Repopurl")
	tpath := os.Getenv("TEMPLATE_PATH")
	if tpath == "" {
		tpath = "./templates"
	}
	log.Println("Loading templates from", tpath)
	err := LoadTemplates(tpath)
	if err != nil {
		panic(err)
	}

	staticFilePath = os.Getenv("STATIC_PATH")
	if staticFilePath == "" {
		staticFilePath = "./static"
	}
	log.Println("Static files at", staticFilePath)

	rootRedirect = os.Getenv("ROOT_REDIRECT")
	log.Println("Redirecting / to", rootRedirect)

	fedoraUsername = os.Getenv("FEDORA_USER")
	fedoraPassword = os.Getenv("FEDORA_PASS")

	dbLocation := os.Getenv("MYSQL_CONNECTION")
	if dbLocation == "" {
		dbLocation = "/repopurl"
		log.Println("Using default database", dbLocation)
	}
	// don't log dblocation otherwise, since it probably contains passwords
	datasource, err = NewMySQL(dbLocation)
	if err != nil {
		log.Println(err)
	}

	router := NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
