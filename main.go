package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	// mySql information for login
	port := os.Getenv("MYSQL_PORT")
	var mysqlLocation string
	if(port == "") {
		mysqlLocation = fmt.Sprintf("%s:%s@%s/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_DB"),
		)
	} else {
		mysqlLocation = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DB"),
		)
	}
	datasource = NewDBSource(mysqlLocation)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

