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

func TestPurlShowFile(t *testing.T) {
	source.db.ExpectBegin()
	source.db.ExpectExec
}

func TestFindPurl(t *testing.T) {
}

func TestCreatePurl(t *testing.T) {
}

func init() {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stud db", err)
	}

	source = &purldb{db: mock}
}
