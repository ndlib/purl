package main

import (
	"fmt"
	"log"
	"net/http"

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

	// config
	var (
		mysqlLocation string
		config Config
	)
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		log.Printf("Error getting config information: %s", err.Error())
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
		source *purldb
		db    *sql.DB
	)
	db, source = NewDBSource(db, mysqlLocation)
	defer db.Close()

	datasource = source

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
