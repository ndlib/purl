package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// A Repository defines the actions we need to do against the
// PURL database. It is an interface so we can use either a MySQL or
// Postgres backend (or an in-memory one for testing).
type Repository interface {
	AllPurls() []Purl

	SummaryStats() Stats

	// FindPurl returns information about the given purl identifier.
	// It returns the zero Purl if there is no purl with that id.
	FindPurl(id int) (Purl, bool)

	FindQuery(query string) []Purl

	LogAccess(vars *http.Request, purl Purl)
}

// A Purl represents a single redirect entry in the database.
// (The database stores this information accross two tables, but we don't need
// to know that here).
type Purl struct {
	ID           int
	RepoID       int
	AccessCount  int
	LastAccessed time.Time
	DateCreated  time.Time
	Filename     string
	URL          string
	Information  string
}

// A RepoObj is a braindead structure needed because we are mirroring
// how the PURLs are stored in the database. For the most part there
// is always a one-to-one relationship between a Purl and a RepoObj.
type RepoObj struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	URL          string    `json:"URL"`
	DateAdded    time.Time `json:"date_added"`
	AddSourceIP  string    `json:"add_source_ip"`
	DateModified time.Time `json:"date_modified"`
	Information  string    `json:"information"`
}

type Stats struct {
	Total         int
	MostUsed      int
	TotalToday    int
	MostUsedToday int
}

var (
	templates            *template.Template
	attachmentExtentions = regexp.MustCompile(`(ovf|zip|vmdk)$`)
	redirectPattern      = regexp.MustCompile(`^(CurateND - |Reformatting Unit:)`)
	fedoraUsername       string
	fedoraPassword       string
	datasource           Repository
	rootRedirect         string
)

// LoadTemplates will load and compile our templates into memory
func LoadTemplates(path string) error {
	var err error
	templates, err = template.ParseGlob(filepath.Join(path, "*"))
	return err
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	err := templates.ExecuteTemplate(w, "404", nil)
	if err != nil {
		log.Println(err)
	}
}

func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	err := templates.ExecuteTemplate(w, "500", nil)
	if err != nil {
		log.Println(err)
	}
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	err := templates.ExecuteTemplate(w, "500", nil)
	if err != nil {
		log.Println(err)
	}
}

// IndexHandler responds to the root route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if rootRedirect != "" {
		http.Redirect(w, r, rootRedirect, 302)
		return
	}
	notFound(w)
}

// AdminHandler returns the number of purls to w.
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	stats := datasource.SummaryStats()
	err := templates.ExecuteTemplate(w, "admin", &stats)
	if err != nil {
		log.Println(err)
	}
}

func AdminSearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("q")
	log.Println("q=", q)
	results := datasource.FindQuery(q)
	err := templates.ExecuteTemplate(w, "admin-search", results)
	if err != nil {
		log.Println(err)
	}
}

// PurlShow returns metadata for the given PURL to w.
func PurlShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purlID, err := strconv.Atoi(vars["purlId"])
	if err != nil {
		notFound(w)
		return
	}
	purl, ok := datasource.FindPurl(purlID)
	if !ok {
		notFound(w)
		return
	}
	err = templates.ExecuteTemplate(w, "view", &purl)
	if err != nil {
		log.Println(err)
	}
}

// PurlShowFile returns either the upstream content of this PURL or
// a redirect to the upstream content. The decision depends on the contents
// of the Information field in the PURL.
func PurlShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purlID, err := strconv.Atoi(vars["purlId"])
	if err != nil {
		notFound(w)
		return
	}

	purl, ok := datasource.FindPurl(purlID)
	if !ok {
		notFound(w)
		return
	}

	// Some entries need a redirect and not a proxy. Determining that from
	// special patterns in the information string is legacy behavior.
	if redirectPattern.MatchString(purl.Information) {
		datasource.LogAccess(r, purl)
		http.Redirect(w, r, purl.URL, 302)
		return
	}

	proxyRequest, _ := http.NewRequest("GET", purl.URL, nil)
	// This test is a little hokey. We assume the fedora URLs are
	// exactly those with the word "fedora" in it somewhere.
	if strings.Contains(purl.URL, "fedora") {
		if fedoraUsername != "" || fedoraPassword != "" {
			proxyRequest.SetBasicAuth(fedoraUsername, fedoraPassword)
		}
	}
	resp, err := http.DefaultClient.Do(proxyRequest)
	if err != nil {
		log.Println(purl.URL, ":", err)
		serverError(w)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("upstream returned", resp.StatusCode, purl.URL)
		serverError(w)
		return
	}

	datasource.LogAccess(r, purl)

	// For certain file extensions, we set the download to be an "attachment"
	// so a web browser will not try to open it in the browser window.
	//
	// We use the filename passed in through the URL, and not the filename
	// stored in the purl record.
	//
	// All of this is previous behavior. Maybe it could be re-examined.
	var disposition string
	filename := purl.Filename
	if attachmentExtentions.MatchString(filename) {
		disposition = "attachment; filename=" + filename
	} else {
		disposition = "inline; filename=" + filename
	}
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	if resp.ContentLength > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(resp.ContentLength, 10))
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println(err)
	}
}

// Helper to set http.ResponseWriter
func setResponseContent(w http.ResponseWriter, r *http.Response, filename string) {
}
