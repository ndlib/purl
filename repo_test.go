// +build mysql

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

	result := datasource.AllPurls()

	for _, res := range result {
		assert.NotEqual(res.Date_created, time.Time{}, "Time incorrectly set on repo")
		assert.NotEqual(res.Id, nil, "Id nil")
	}
}

func TestFindPurl(t *testing.T) {
	assert := assert.New(t)

	result := datasource.FindPurl(5)

	assert.NotEqual(result.Date_created, time.Time{}, "Time incorrectly set on repo")
	assert.NotEqual(result.Id, nil, "Id nil")

	assert.Equal(result.Id, 5, "Id not correct")
	assert.Equal(result.Repo_obj_id, "5", "Repo Id not correct")
	assert.Equal(result.Access_count, 625, "Id not correct")
	time_val, _ := time.Parse(time.RFC3339, "2016-11-15T14:16:14Z")
	assert.Equal(result.Last_accessed, time_val, "Last_accesed not correct")
	time_val, _ = time.Parse(time.RFC3339, "2011-09-14T14:40:11Z")
	assert.Equal(result.Date_created, time_val, "Date_created not correct")
}

func TestCreatePurl(t *testing.T) {
	assert := assert.New(t)

	var newpurl = Purl{
		Id:          11,
		Repo_obj_id: "110",
	}
	newpurl.Last_accessed, _ = time.Parse(time.RFC3339, "2016-11-16T03:33:33Z")
	newpurl.Date_created, _ = time.Parse(time.RFC3339, "2011-09-14T13:55:55Z")

	datasource.CreatePurl(newpurl)

	result := datasource.FindPurl(11)

	assert.Equal(result.Id, newpurl.Id, "Id not correct")
	assert.Equal(result.Repo_obj_id, newpurl.Repo_obj_id, "Repo Id not correct")
	assert.Equal(result.Last_accessed, newpurl.Last_accessed, "Last_accesed not correct")
	assert.Equal(result.Date_created, newpurl.Date_created, "Date_created not correct")

	// _ = datasource.destroyPurlDB(11)
}

func init() {
	connection := os.Getenv("MYSQL_CONNECTION")
	if connection == "" {
		connection = "/test"
		fmt.Println("MYSQL_CONNECTION not set. Using default:", connection)
	}

	datasource = NewDBSource(connection)
}
