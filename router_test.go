// +build router

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var (
	source *memoryRepo
)

func TestRouterProxy(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)
	server := httptest.NewServer(router)
	defer server.Close()

	newpurl := Purl{
		Id:          11,
		Repo_obj_id: "110",
	}

	newrepo := RepoObj{
		Id:          110,
		Information: `CurateND - |Reformatting Unit:`,
		Url:         `http://catalog.hathitrust.org/Record/009783954`,
	}

	datasource.CreatePurl(newpurl)
	datasource.CreateRepo(newrepo)

	req, err := http.NewRequest("GET", `/view/11/any.pdf`, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	NewRouter().ServeHTTP(rr, req)
	if rr.Code != 302 {
		t.Errorf("invalid status code", rr.Code)
	}
}

func init() {
	repo := &memoryRepo{}
	datasource = repo
}
