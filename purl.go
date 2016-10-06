package main

import (
	"fmt"
	"io/ioutil"
  "encoding/json"
	"log"
	"net/http"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type PURL struct {
  purl_id int
  file_name string
	repo_url string
  purl_url string
  access_count int
  last_accessed int
	description string
}

func (p *URL) save() error {
  fileurl := p.local_url
}

func main() {
  db, err := sql.Open("mysql",
                    "root@tcp(127.0.0.1:3306)")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  if db.Ping() {
    fmt.Println("good")
  } else{
    log.Fatal(err)
  }

	// return query from database
	rows, err := db.Query("select * from purl where purl_id = ?", 1)
	if err != nill{
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		err:= rows.Scan(&purl_id,&file_name)
		if err!= nil {
			log.Fatal(err)
		}
		log.Println(purl_id, file_name)
	}
	err = rows.Err(
		if err != nill{
			log.Fatal(err)
		}
	)
}
