package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var repoRoutes = []route{
	{"GET", "/", Index},
	{"GET", "/admin", AdminIndex},
	{"GET", "/purls", PurlIndex},
	{"POST", "/purl/create", PurlCreate},
	{"GET", "/view/{purlId}", PurlShow},
	{"GET", "/view/{purlId}/{filename}", PurlShowFile},
	{"GET", "/query?={query}", Query},
}

// NewRouter returns a Handler that will take care of all the repopurl routes.
func NewRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range repoRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandlerFunc)
	}

	return &LogHandler{router}
}

// LogHandler wraps a handler and logs the request path and duration at the end.
type LogHandler struct {
	h http.Handler
}

func (lh *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lh.h.ServeHTTP(w, r)
	log.Println(r.Method, r.RequestURI, time.Since(start))
}
