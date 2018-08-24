// +build go1.8

// The proxy redirect test relies on the ErrUseLastResponse code
// which first appeared in go 1.8

package main

import (
	"net/http"
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
		resp.Body.Close()
	}
}
