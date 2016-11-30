package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gcfg.v1"
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

	// repo := &memoryRepo{}
	// Give us some seed data
	// repo.CreatePurl(Purl{source_app: "Library"})
	// repo.CreatePurl(Purl{source_app: "Code School"})

	// config
	var (
		// port          string
		// storageDir    string
		// logfilename   string
		// sqliteFile    string
		mysqlLocation string
		// showVersion   bool
		// configFile    string
		config Config
	)
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		panic(err)
	}

	// mySql information for login
	mysqlLocation = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Host,
		config.Mysql.Port,
		config.Mysql.Database,
	)

	// Connect to mysql database
	var (
		store *purldb
		db    *sql.DB
		// errDB error
	)
	db, err = sql.Open("mysql", mysqlLocation)
	fmt.Println(mysqlLocation)
	if err != nil {
		log.Printf("Error opening database: %s", err.Error())
	}
	defer db.Close()
	if db != nil {
		var wait = 1
		store = NewDBFileStore(db)
		for store == nil {
			log.Printf("Problem loading pools from database. Trying again in %d seconds", wait)
			time.Sleep(time.Duration(wait) * time.Second)
			wait *= 2
			if wait > 300 {
				wait = 300
			}
			store = NewDBFileStore(db)
		}
	}
	// SetupHandlers(store)
	// log.Println("Listening on port", port)
	// err = http.ListenAndServe(":"+port, nil)
	// if err != nil {log.Fatal("ListenAndServe", err)}
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %s", err.Error())
	}
	datasource = store

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
