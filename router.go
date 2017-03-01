package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var repoRoutes = routes{
	route{
		"Index",
		"GET",
		"/",
		Index,
	},
	route{
		"AdminIndex",
		"GET",
		"/admin",
		AdminIndex,
	},
	route{
		"PurlIndex",
		"GET",
		"/purls",
		PurlIndex,
	},
	route{
		"PurlCreate",
		"POST",
		"/purl/create",
		PurlCreate,
	},
	route{
		"PurlShow",
		"GET",
		"/view/{purlId}",
		PurlShow,
	},
	route{
		"PurlShowFile",
		"GET",
		"/view/{purlId}/{filename}",
		PurlShowFile,
	},
	route{
		"Query",
		"GET",
		"/query?={query}",
		Query,
	},
}

// Our initial router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range repoRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
