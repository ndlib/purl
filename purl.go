package main

import (
	"log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type PURL struct {
  purl_id int
  note string
  file_name string
  last_accessed string
  local_url string
  redirect_uri string
}

// func (p *URL) save() error {
//   fileurl := p.local_url
// }

func main() {
  db, err := sql.Open("mysql",
                    "root:Yonzalk12@tcp(127.0.0.1:3306)/purl")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  err = db.Ping()
	if err != nil {
    log.Fatal(err)
  }
}
