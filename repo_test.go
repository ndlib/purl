// +build mysql

// These tests assume the test/seed_data.sql file has been loaded
// into the MySQL database.

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllPurls(t *testing.T) {
	assert := assert.New(t)

	result := mysqlTarget.AllPurls()

	for _, res := range result {
		t.Log(res)
		assert.NotEqual(res.DateCreated, time.Time{}, "Time incorrectly set on repo")
		assert.NotEqual(res.ID, nil, "ID nil")
	}
}

func TestFindPurl(t *testing.T) {
	assert := assert.New(t)

	result := mysqlTarget.FindPurl(5)

	assert.NotEqual(result.DateCreated, time.Time{}, "Time incorrectly set on repo")
	assert.NotEqual(result.ID, nil, "ID nil")

	assert.Equal(result.ID, 5, "ID not correct")
	assert.Equal(result.RepoObjID, "5", "Repo ID not correct")
	assert.Equal(result.AccessCount, 625, "AccessCount not correct")
	time_val, _ := time.Parse(time.RFC3339, "2016-11-15T14:16:14Z")
	assert.Equal(result.LastAccessed, time_val, "LastAccesed not correct")
	time_val, _ = time.Parse(time.RFC3339, "2011-09-14T14:40:11Z")
	assert.Equal(result.DateCreated, time_val, "DateCreated not correct")
}

func TestCreatePurl(t *testing.T) {
	assert := assert.New(t)

	var newpurl = Purl{
		ID:        11,
		RepoObjID: "110",
	}
	newpurl.LastAccessed, _ = time.Parse(time.RFC3339, "2016-11-16T03:33:33Z")
	newpurl.DateCreated, _ = time.Parse(time.RFC3339, "2011-09-14T13:55:55Z")

	_ = mysqlTarget.destroyPurlDB(11)
	mysqlTarget.CreatePurl(newpurl)

	result := mysqlTarget.FindPurl(11)

	assert.Equal(result.ID, newpurl.ID, "ID not correct")
	assert.Equal(result.RepoObjID, newpurl.RepoObjID, "Repo ID not correct")
	assert.Equal(result.LastAccessed, newpurl.LastAccessed, "LastAccesed not correct")
	assert.Equal(result.DateCreated, newpurl.DateCreated, "DateCreated not correct")

	_ = mysqlTarget.destroyPurlDB(11)
}

var (
	mysqlTarget *purldb
)

func init() {
	connection := os.Getenv("MYSQL_CONNECTION")
	if connection == "" {
		connection = "/test"
		fmt.Println("MYSQL_CONNECTION not set. Using default:", connection)
	}

	mysqlTarget = NewDBSource(connection)
}
