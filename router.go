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
	{"GET", "/", IndexHandler},
	{"GET", "/admin", AdminHandler},
	{"GET", "/view/{purlId}", PurlShow},
	{"GET", "/view/{purlId}/{filename}", PurlShowFile},
	{"HEAD", "/view/{purlId}/{filename}", PurlShowFile},
	{"GET", "/admin/search", AdminSearchHandler},
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

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFound(w)
	})

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
