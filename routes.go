package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"PurlIndex",
		"GET",
		"/purls",
		PurlIndex,
	},
	Route{
		"PurlCreate",
		"POST",
		"/purls",
		PurlCreate,
	},
	Route{
		"PurlShow",
		"GET",
		"/purls/{purlId}",
		PurlShow,
	},
}
