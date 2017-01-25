// +build mysql

package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	source *purldb
)

// UNAPPROVED: commands to test the sql
// without needing to run a backend 
// database. Implemented but not used
// func TestcreatePurlDB(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' opening db", err)
// 	}
// 	defer db.Close()
// 	source = &purldb{db: db}
// 	purl_test := Purl{
// 		Repo_obj_id:1, 
// 		Access_count:723, 
// 		Last_accessed:"2016-11-16 03:27:15", 
// 		Date_created: "2011-09-14 13:56:38",
// 	}

// 	mock.ExpectBegin()
// 	mock.ExpectExec("INSERT INTO purl (purl_id, repo_object_id, last_accessed, source_app, date_created)").
// 		WithArgs(purl_test).
// 		WillReturnResult(sqlmock.NewResult(1,1))
// 	mock.ExpectCommit()

// 	if err := source.createPurlDB(purl_test); err != nil {
// 		t.Errorf("error there was an unpected failure when testing insert: %s", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectation: %s", err)
// 	}
// }

func TestAllPurls(t *testing.T) {
	assert := assert.New(t)

	var purltestdb *purldb
	purltestdb = source

	result := purltestdb.AllPurls()

	for _, res := range result {
		assert.NotEqual(res.Date_created, time.Time{}, "Time incorrectly set on repo")
		assert.NotEqual(res.Id, nil, "Id nil")
	}
}

func TestFindPurl(t *testing.T) {
	assert := assert.New(t)

	var purltestdb *purldb
	purltestdb = source
	result := purltestdb.FindPurl(5)

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
	purltestdb = source
	_ = purltestdb.createPurlDB(newpurl)

	result := purltestdb.FindPurl(11)

	assert.Equal(result.Id, newpurl.Id, "Id not correct")
	assert.Equal(result.Repo_obj_id, newpurl.Repo_obj_id, "Repo Id not correct")
	assert.Equal(result.Last_accessed, newpurl.Last_accessed, "Last_accesed not correct")
	assert.Equal(result.Date_created, newpurl.Date_created, "Date_created not correct")

	_ = purltestdb.destroyPurlDB(11)
}

func init() {
	mysqlconn := os.Getenv("MYSQL_CONNECTION")
	if mysqlconn == "" {
		panic("MYSQL_CONNECTION not set")
	}
	source = NewDBSource(mysqlconn)
}
