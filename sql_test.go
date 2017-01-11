package main

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mysqldb *sql.DB

var (
	datatest Repository
)

func TestAllPurls(t *testing.T) {
	assert := assert.New(t)

	var purltestdb *purldb
	purltestdb = NewDBFileStore(mysqldb)

	result := purltestdb.AllPurls()

	// assert.Equal(len(result), 10, "We should have ten repos in the database")
	for _, res := range result {
		assert.NotEqual(res.Date_created, time.Time{}, "Time incorrectly set on repo")
		assert.NotEqual(res.Id, nil, "Id nil")
	}
}

func TestFindPurl(t *testing.T) {
	assert := assert.New(t)

	var purltestdb *purldb
	purltestdb = NewDBFileStore(mysqldb)
	result := purltestdb.FindPurl(5)

	// assert.Equal(len(result), 10, "We should have ten repos in the database")
	assert.NotEqual(result.Date_created, time.Time{}, "Time incorrectly set on repo")
	assert.NotEqual(result.Id, nil, "Id nil")

	assert.Equal(result.Id, 5, "Id not correct")
	assert.Equal(result.Repo_obj_id, "5", "Repo Id not correct")
	assert.Equal(result.Access_count, 625, "Id not correct")
	time_val, err := time.Parse(time.RFC3339, "2016-11-15T14:16:14Z")
	if err != nil {
		panic(err)
	}
	assert.Equal(result.Last_accessed, time_val, "Last_accesed not correct")
	time_val, err = time.Parse(time.RFC3339, "2011-09-14T14:40:11Z")
	if err != nil {
		panic(err)
	}
	assert.Equal(result.Date_created, time_val, "Date_created not correct")
}

func TestCreatePurl(t *testing.T) {
	assert := assert.New(t)

	var newpurl Purl
	newpurl.Id = 11
	newpurl.Repo_obj_id = "110"
	var err error
	newpurl.Last_accessed, err = time.Parse(time.RFC3339, "2016-11-16T03:33:33Z")
	if err != nil {
		panic(err)
	}
	newpurl.Date_created, err = time.Parse(time.RFC3339, "2011-09-14T13:55:55Z")
	if err != nil {
		panic(err)
	}

	var purltestdb *purldb
	purltestdb = NewDBFileStore(mysqldb)
	_ = purltestdb.createPurlDB(newpurl)

	result := purltestdb.FindPurl(11)

	assert.Equal(result.Id, newpurl.Id, "Id not correct")
	assert.Equal(result.Repo_obj_id, newpurl.Repo_obj_id, "Repo Id not correct")
	assert.Equal(result.Last_accessed, newpurl.Last_accessed, "Last_accesed not correct")
	assert.Equal(result.Date_created, newpurl.Date_created, "Date_created not correct")

	_ = purltestdb.destroyPurlDB(11)
}

func init() {
	// config
	var (
		mysqlconn string
	)
	// err := gcfg.ReadFileInto(&config, "config.gcfg")
	// if err != nil {
	// 	panic(err)
	// }
	// mysqlconn = fmt.Sprintf("travis@tcp(127.0.0.1:3600)/test_database?parseTime=true")
	mysqlconn = os.Getenv("MYSQL_CONNECTION")
	if mysqlconn == "" {
		panic("MYSQL_CONNECTION not set")
	}
	var err_db error
	mysqldb, err_db = sql.Open("mysql", mysqlconn+"?parseTime=true")
	if err_db != nil {
		panic(err_db)
	}
	err_ping := mysqldb.Ping()
	if err_ping != nil {
		log.Printf("Error pinging database: %s", err_ping.Error())
	}
	// defer mysqldb.Close()
}
