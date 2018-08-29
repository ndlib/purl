// repo-munge is a simple command line utility to either
// munge the URLs or unmunge them. A munging simply appends
// a hyphen to the end. Use the environment variable
// MYSQL_CONNECTION to set the database to work against.
//
// usage:
//
//     export "MYSQL_CONNECTION=user:pass@tcp(hostname)/repopurl"
//     ./repo-munge
//     ./repo-munge -unmunge
//
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	unmunge := flag.Bool("unmunge", false, "Undo the munging")
	flag.Parse()

	location := os.Getenv("MYSQL_CONNECTION")
	db, err := sql.Open("mysql", location)
	if err != nil {
		log.Println(err)
		return
	}
	var command string
	if *unmunge {
		log.Println("Unmunging rows")
		command = `UPDATE repo_object
		SET url=regexp_replace(url, "-+$", "")
		WHERE information LIKE "Catholic Pamphlet%"`
	} else {
		log.Println("Munging rows")
		command = `UPDATE repo_object
		SET url=concat(url, "-")
		WHERE information LIKE "Catholic Pamphlet%"`
	}

	result, err := db.Exec(command)
	if err != nil {
		log.Println(err)
	}
	nrows, _ := result.RowsAffected()
	log.Println(nrows, "rows updated")
	db.Close()
}
