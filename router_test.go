package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestShowFile(t *testing.T) {
	// are the requests that should be proxied proxied?
	table := []URLTest{
		{path: "/view/abcdefg/any.pdf", status: 404},
		{path: "/view/-1/any.pdf", status: 404},
		{path: "/view/123456789/any.pdf", status: 404},
		{path: "/view/500/any.pdf", status: 200, body: "a very good file"},
		// upstream is 404
		{path: "/view/501/any.pdf", status: 500},
		// redirects? (vs proxy)
		{path: "/view/502/any.pdf", status: 500, body: "hello world"},
		{path: "/view/503/any.pdf", status: 200, body: "hello world"},
		// propagate content type from upstream?
		{path: "/view/504/any.pdf", status: 200, body: "hello world",
			headers: map[string]string{
				"Content-Type":        "application/qqq",
				"Content-Disposition": "inline; filename=best.pdf"}},
		// propagage content length from upstream?
		{path: "/view/505/any.pdf", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "inline; filename=longfilename.pdf"}},
		// test the inline/attachment switch (.zip, .ovf, .vmdk extensions)
		{path: "/view/505/any.zip", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "inline; filename=longfilename.pdf"}},
		{path: "/view/505/any.ovf", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "inline; filename=longfilename.pdf"}},
		{path: "/view/505/any.vmdk", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "inline; filename=longfilename.pdf"}},
		{path: "/view/506/any.pdf", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "attachment; filename=longfilename.vmdk"}},
		{path: "/view/507/any.pdf", status: 200, body: "a very",
			headers: map[string]string{
				"Content-Length":      "6",
				"Content-Disposition": "attachment; filename=longfilename.zip"}},
		// suppressed record
		{path: "/view/508/download", status: 404},
	}

	for _, test := range table {
		checkSimpleGetRequest(t, repoServer, test)
	}
}

type URLTest struct {
	path    string
	status  int
	body    string
	headers map[string]string
}

// checkSimpleGetRequest does a GET request to URL, and then compares the response code and
// response body to what was provided. Any errors are flagged on t.
func checkSimpleGetRequest(t *testing.T, server *httptest.Server, test URLTest) {
	URL := server.URL + test.path
	resp, err := http.Get(URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != test.status {
		t.Errorf("On %s received status %d, expected %d", test.path, resp.StatusCode, test.status)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if test.body != "" && string(b) != test.body {
		t.Errorf("On %s received body: %s\n    expected: %s", test.path, b, test.body)
	}
	//t.Log(resp)
	for header, expected := range test.headers {
		received := resp.Header.Get(header)
		if received != expected {
			t.Errorf("On %v got %v: %v\n    expected %v", test.path, header, received, expected)
		}
	}
	resp.Body.Close()
}

func TestLogAccess(t *testing.T) {
	tests := []struct {
		ID        int
		path      string
		status    int
		increment bool
	}{
		{500, "/view/500/a", 200, true},  // usual case
		{500, "/view/500", 200, false},   // record view page
		{501, "/view/501/a", 500, false}, // error on proxy
		{502, "/view/502/a", 500, true},  // redirect to error
		{503, "/view/503/a", 200, true},  // redirect to good
	}
	for _, test := range tests {
		purlv1, ok := datasource.FindPurl(test.ID)
		if !ok {
			t.Error("Couldn't find purl", test.ID)
			continue
		}
		checkSimpleGetRequest(t, repoServer, URLTest{path: test.path, status: test.status})
		purlv2, _ := datasource.FindPurl(test.ID)
		delta := purlv2.AccessCount - purlv1.AccessCount
		if test.increment && delta <= 0 {
			t.Error("Found count", purlv2.AccessCount, "expected something larger than", purlv1.AccessCount)
		} else if !test.increment && delta > 0 {
			t.Error("Found count", purlv2.AccessCount, "expected to be equal to", purlv1.AccessCount)
		}
	}
}

var (
	repoServer  *httptest.Server
	dummyServer *httptest.Server
)

func init() {
	LoadTemplates("./templates")
	memory := &memoryRepo{}
	// have the handlers reference our test store
	datasource = memory

	// set up the repo server AND ALSO a second dummy server that will be a proxy source.
	repoServer = httptest.NewServer(NewRouter())
	dummyServer = httptest.NewServer(http.HandlerFunc(dummyHandler))

	// now seed data that points to the dummy server
	seedItems := []Purl{
		{
			ID:          500,
			Filename:    "good.pdf",
			URL:         dummyServer.URL + "/200?data=a+very+good+file",
			Information: "",
		},
		{
			ID:          501,
			Filename:    "bad.pdf",
			URL:         dummyServer.URL + "/404",
			Information: "item title",
		},
		{
			ID:          502,
			Filename:    "redirect",
			URL:         dummyServer.URL + "/500",
			Information: "CurateND - item page",
		},
		{
			ID:          503,
			Filename:    "redirect",
			URL:         dummyServer.URL + "/200",
			Information: "Reformatting Unit: item name",
		},
		{
			ID:          504,
			Filename:    "best.pdf",
			URL:         dummyServer.URL + "/200?type=application/qqq",
			Information: "",
		},
		{
			ID:          505,
			Filename:    "longfilename.pdf",
			URL:         dummyServer.URL + "/200?data=a+very+long+text&size=6",
			Information: "",
		},
		{
			ID:          506,
			Filename:    "longfilename.vmdk",
			URL:         dummyServer.URL + "/200?data=a+very+long+text&size=6",
			Information: "",
		},
		{
			ID:          507,
			Filename:    "longfilename.zip",
			URL:         dummyServer.URL + "/200?data=a+very+long+text&size=6",
			Information: "",
		},
		{
			ID:          508,
			Filename:    "suppressed",
			URL:         dummyServer.URL + "/200-",
			Information: "",
		},
	}
	for _, seed := range seedItems {
		seed.RepoID = seed.ID
		memory.CreatePurl(seed)
	}
}

// dummyHandler is for testing. The path is of the form /{status code}.
// The "data" parameter can pass the data to be returned in the body; it
// defaults to "hello world". The "type" parameter will set the return
// content type; defaults to sniffing the data. The "size" patameter will set the
// return content length. defaults to the length of the data.
//
// added the size parameter to test for negative content lengths (fedora HEAD
// request bug). but the go library really hates negative lengths and does
// its best to be correct, making this not really useful to test. Maybe remove
// the functionality?
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	// remove initial "/"
	status, _ := strconv.Atoi(r.URL.Path[1:])
	data := r.FormValue("data")
	typ := r.FormValue("type")
	size := r.FormValue("size")
	if status <= 0 || status >= 1000 {
		// REALLY bad status. normalize it
		status = 400
	}
	if data == "" {
		data = "hello world"
	}
	if typ != "" {
		w.Header().Set("Content-Type", typ)
	}
	if size != "" {
		w.Header().Set("Content-Length", size)
		m, _ := strconv.Atoi(size)
		if m < len(data) {
			data = data[:m]
		}
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", data)
}

// of course we test our dummy handler!
func TestDummyHandler(t *testing.T) {
	table := []URLTest{
		{path: "/200", status: 200, body: "hello world"},
		{path: "/200?data=a+fine+day", status: 200, body: "a fine day"},
		{path: "/333", status: 333, body: "hello world"},
		{path: "/404?data=a", status: 404, body: "a"},
		{path: "/200?type=foo/bar", status: 200, body: "hello world",
			headers: map[string]string{"Content-Type": "foo/bar"}},
		// the server sends all the data, but we only read 0 bytes of it
		{path: "/200?size=0", status: 200, body: "",
			headers: map[string]string{"Content-Length": "0"}},
		{path: "/200?data=abcdefgh&size=3", status: 200, body: "abc"},
	}

	for _, test := range table {
		checkSimpleGetRequest(t, dummyServer, test)
	}
}
