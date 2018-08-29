// +build mysql

// These tests assume the test/seed_data.sql file has been loaded
// into the MySQL database.

package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMySQLAllPurls(t *testing.T) {
	assert := assert.New(t)

	result := mysqlTarget.AllPurls()

	for _, res := range result {
		t.Log(res)
		assert.NotEqual(res.DateCreated, time.Time{}, "Time incorrectly set on repo")
		assert.NotEqual(res.ID, nil, "ID nil")
	}
}

func TestMySQLFindPurl(t *testing.T) {
	assert := assert.New(t)

	result, ok := mysqlTarget.FindPurl(5)
	if !ok {
		t.Fatal("Couldn't find record")
	}

	assert.NotEqual(result.DateCreated, time.Time{}, "Time incorrectly set on repo")
	assert.NotEqual(result.ID, nil, "ID nil")

	assert.Equal(5, result.ID, "ID not correct")
	assert.Equal(5, result.RepoID, "Repo ID not correct")
	assert.Equal(625, result.AccessCount, "AccessCount not correct")
	tv, _ := time.Parse(time.RFC3339, "2016-11-15T14:16:14Z")
	assert.Equal(tv, result.LastAccessed, "LastAccesed not correct")
	tv, _ = time.Parse(time.RFC3339, "2011-09-14T14:40:11Z")
	assert.Equal(tv, result.DateCreated, "DateCreated not correct")
}

func getaccesscount(purlID int) int {
	var count int
	mysqlTarget.db.QueryRow("SELECT access_count FROM purl WHERE purl_id = ?", purlID).Scan(&count)
	return count
}

func TestMySQLFindQuery(t *testing.T) {
	tests := []struct {
		query string
		nrows int
	}{
		{"", 10},
		{"Catholic", 4},
		{"catholic", 4},
		{"papers", 0},
		{"CurateND", 1},
	}

	for _, test := range tests {
		result := mysqlTarget.FindQuery(test.query)
		if len(result) != test.nrows {
			t.Errorf("For '%s' expected %d rows, received %d\n",
				test.query, test.nrows, len(result))
		}
	}
}

func TestMySQLLogAccess(t *testing.T) {
	// get starting count, date/time
	firstCount := getaccesscount(10)

	// does this update last accessed?
	purl, _ := mysqlTarget.FindPurl(10)

	secondCount := getaccesscount(10)
	if firstCount != secondCount {
		t.Error("access count changed", firstCount, secondCount)
	}

	// does this update last accessed?
	req, _ := http.NewRequest("GET", "/", nil)
	mysqlTarget.LogAccess(req, purl)

	thirdCount := getaccesscount(10)
	if thirdCount != firstCount+1 {
		t.Error("Found", thirdCount, "expected", firstCount+1)
	}
}

var (
	mysqlTarget *mysqlDB
)

func init() {
	connection := os.Getenv("MYSQL_CONNECTION")
	if connection == "" {
		connection = "/test"
		fmt.Println("MYSQL_CONNECTION not set. Using default:", connection)
	}

	r, err := NewMySQL(connection)
	if err != nil {
		panic(err)
	}
	mysqlTarget = r.(*mysqlDB)
}
