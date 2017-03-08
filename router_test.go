package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRouterProxy(t *testing.T) {
	// are the requests that should be proxied proxied?
	table := []struct {
		path   string
		status int
	}{
		{path: "/view/502/any.pdf", status: 302},
		{path: "/view/503/any.pdf", status: 302},
		{path: "/view/500/any.pdf", status: 200},
		{path: "/view/501/any.pdf", status: 500},
	}

	// use custom client so we don't follow redirects (since we want to TEST
	// whether redirects were returned!)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, test := range table {
		req, err := http.NewRequest("GET", repoServer.URL+test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if resp.StatusCode != test.status {
			t.Errorf("On %s received status %d, expected %d", test.path, resp.StatusCode, test.status)
		}
		b, _ := ioutil.ReadAll(resp.Body)
		t.Logf("On %s, body: %s", test.path, b)
		resp.Body.Close()
	}
}

func TestShowFile(t *testing.T) {
	// are the requests that should be proxied proxied?
	table := []URLTest{
		{path: "/view/abcdefg/any.pdf", status: 404, body: `{"code":404,"text":"Not Found"}` + "\n"},
		{path: "/view/-1/any.pdf", status: 404, body: `{"code":404,"text":"Not Found"}` + "\n"},
		{path: "/view/123456789/any.pdf", status: 404, body: `{"code":404,"text":"Not Found"}` + "\n"},
		{path: "/view/500/any.pdf", status: 200, body: "a very good file"},
		{path: "/view/501/any.pdf", status: 500, body: "Content Unavailable\n"},
		{path: "/view/502/any.pdf", status: 500, body: "hello world"},
		{path: "/view/503/any.pdf", status: 200, body: "hello world"},
	}

	for _, test := range table {
		checkSimpleGetRequest(t, repoServer.URL+test.path, test)
	}
}

type URLTest struct {
	path   string
	status int
	body   string
}

// checkSimpleGetRequest does a GET request to URL, and then compares the response code and
// response body to what was provided. Any errors are flagged on t.
func checkSimpleGetRequest(t *testing.T, URL string, test URLTest) {
	resp, err := http.Get(URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != test.status {
		t.Errorf("On %s received status %d, expected %d", test.path, resp.StatusCode, test.status)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if string(b) != test.body {
		t.Errorf("On %s received body: %s\n    expected: %s", test.path, b, test.body)
	}
	resp.Body.Close()
}

var (
	repoServer  *httptest.Server
	dummyServer *httptest.Server
)

func init() {
	memory := &memoryRepo{}
	// have the handlers reference our test store
	datasource = memory

	// set up the repo server AND ALSO a second dummy server that will be a proxy source.
	repoServer = httptest.NewServer(NewRouter())
	dummyServer = httptest.NewServer(http.HandlerFunc(dummyHandler))

	// now seed data that points to the dummy server
	seedItems := []RepoObj{
		{
			Id:          500,
			Filename:    "good.pdf",
			Url:         dummyServer.URL + "/200?data=a+very+good+file",
			Information: "",
		},
		{
			Id:          501,
			Filename:    "bad.pdf",
			Url:         dummyServer.URL + "/404",
			Information: "item title",
		},
		{
			Id:          502,
			Filename:    "redirect",
			Url:         dummyServer.URL + "/500",
			Information: "CurateND - item page",
		},
		{
			Id:          503,
			Filename:    "redirect",
			Url:         dummyServer.URL + "/200",
			Information: "Reformatting Unit: item name",
		},
	}
	for _, seed := range seedItems {
		memory.CreateRepo(seed)
		memory.CreatePurl(Purl{
			Id:          seed.Id,
			Repo_obj_id: fmt.Sprintf("%d", seed.Id),
		})
	}
}

// dummyHandler is for testing. The path is of the form /{status code}.
// The "data" parameter can pass the data to be returned in the body; it
// defaults to "hello world".
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	// remove initial "/"
	status, _ := strconv.Atoi(r.URL.Path[1:])
	data := r.FormValue("data")
	if status < 0 || status >= 1000 {
		// REALLY bad status. normalize it
		status = 400
	}
	if data == "" {
		data = "hello world"
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, data)
}

// of course we test our dummy handler!
func TestDummyHandler(t *testing.T) {
	table := []URLTest{
		{path: "/200", status: 200, body: "hello world"},
		{path: "/200?data=a+fine+day", status: 200, body: "a fine day"},
		{path: "/333", status: 333, body: "hello world"},
		{path: "/404?data=a", status: 404, body: "a"},
	}

	for _, test := range table {
		checkSimpleGetRequest(t, dummyServer.URL+test.path, test)
	}
}
