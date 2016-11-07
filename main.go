package main

import (
	"log"
	"net/http"
	"time"
)

var (
	datasource Repository
)

func main() {

	repo := &memoryRepo{}
	// Give us some seed data
	repo.CreatePurl(Purl{Id:45,
												Access_count:3,
		 										Last_accessed:time.Now(),
												Repo_url:"http://fedoraprod.library.nd.edu:8080/fedora/get/CATHOLLIC-PAMPHLET:743445/PDF1",
												Purl_url:"http://repopurlpprd.library.nd.edu/view/45/743445.pdf",
												Repo_obj_id:"743445",
											File_name:"743445.pdf"})

	datasource = repo

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
