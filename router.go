package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var repoRoutes = []route{
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

// Logger returns a Handler that wraps inner and logs the request path and duration at the end.
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
