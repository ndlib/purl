// +build proxy

package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_  "github.com/DATA-DOG/go-sqlmock"
)

var (
	source *purldb
)

func TestAllPurls(t *testing.T) {
}

func TestFindPurl(t *testing.T) {
}

func TestCreatePurl(t *testing.T) {
}

func init() {
	mysqlconn := os.Getenv("MYSQL_CONNECTION")
	if mysqlconn == "" {
		panic("MYSQL_CONNECTION not set")
	}
	source = NewDBSource(mysqlconn)
}
