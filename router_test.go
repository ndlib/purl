// +build router

package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	source memoryRepo
)

func TestRouterProxy(t *testing.T) {
	server := httptest.NewServer(NewRouter())
	defer server.Close()

	var newpurl Purl
	var err error
	newpurl.Id = 11
	newpurl.Repo_obj_id = "110"
	newpurl.Last_accessed, err = time.Parse(time.RFC3339, "2016-11-16T03:33:33Z")
	if err != nil {
		t.Error(err.Error())
	}
	newpurl.Date_created, err = time.Parse(time.RFC3339, "2011-09-14T13:55:55Z")
	if err != nil {
		panic(err)
	}

	var newrepo RepoObj
	newrepo.Id = 110
	newrepo.Information = `^(CurateND - |Reformatting Unit:)`
	newrepo.Url = `www.example.com`

	source.CreatePurl(newpurl)
	source.CreateRepo(newrepo)

	res, err := httptest.Get(server.URL + `/view/11/any.pdf`)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 302 {
		t.Error("invalid status code")
	}
}

func init() {
	repo := &memoryRepo{}
	source = repo
}
