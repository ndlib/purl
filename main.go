package main

import (
	"log"
	"net/http"

	"database/sql"
	"gopkg.in/gcfg.v1"
	_ "github.com/go-sql-driver/mysql"
)

var (
	datasource Repository
)

type Config struct {
	General struct {
		Port       string
		StorageDir string
	}
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     string
		Database string
	}
}

func main() {

	repo := &memoryRepo{}
	// Give us some seed data
	// repo.CreatePurl(Purl{source_app: "Library"})
	// repo.CreatePurl(Purl{source_app: "Code School"})

	// config
	var config Config
	err := gcfg.ReadFileInto(&config, "config.gcfg")

	// mySql information for login
	mysqlLocation = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				config.Mysql.User,
				config.Mysql.Password,
				config.Mysql.Host,
				config.Mysql.Port,
				config.Mysql.Database
	)

	// Connect to mysql database
	var (
		store purl_db
		db *sql.DB
		err error
	)
	db, err = sql.Open("mysql", mysqlLoc)
	if err != nil(
		log.Fatalf("Error opening database: %s", err.Error())
	)
	if db != nil {
		var wait = 1
		store = NewDBFileStore(db)
		while store == nil {
			time.Sleep(time.Duration(wait)*time.Second)
			wait *= 2
			if wait > 300 {
				wait = 300
			}
			store = NewDBFileStore(db)
		}
	}

	datasource = repo

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
