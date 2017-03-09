package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var regexpCurate *regexp.Regexp = regexp.MustCompile(`^(CurateND - |Reformatting Unit:)`)
var reHttp *regexp.Regexp = regexp.MustCompile(`http(s?)://(.+)`)
var reZip *regexp.Regexp = regexp.MustCompile(`\b(ovf$)|\b(zip$)|\b(vmdk$)`)

var (
	txViewTemplate = template.Must(template.New("txinfo").Parse(`<html>
	<head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	    <title>Hesburgh Libraries Permanent URL System</title>
	</head>

	<body lang="en">

	    <h1>View PURL</h1>

	    <table>
	        <tbody>
	            <tr>
	                <td>ID</td>
	                <td>{{.Id}}</td>
	            </tr>
	            <tr>
	                <td>Note</td>
	                <td>{{.Information}}</td>
	            </tr>
	            <tr>
	                <td>File Name</td>
	                <td>{{.File_name}}</td>
	            </tr>
	            <tr>
	                <td>Last Accessed</td>
	                <td>{{.Last_accessed}}</td>
	            </tr>
	            <tr>
	                <td>Repository URL</td>
	                <td><a href="{{.Repo_url}}">{{.Repo_url}} </a></td>
	            </tr>
	            <tr>
	                <td>Access Count</td>
	                <td>{{.Access_count}}</td>
	            </tr>
	        </tbody>
	    </table>

	    <a title="University of Notre Dame" href="http://www.nd.edu/">University of Notre Dame</a>
	</body></html>`))
)

// sendJSON set the response code to the given status, and then writes the JSON
// serialization of data to w. Errors are logged.
func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func sendNotFound(w http.ResponseWriter) {
	data := jsonErr{Code: http.StatusNotFound, Text: "Not Found"}
	sendJSON(w, data, http.StatusNotFound)
}

// HELPERS FOR THE HANDLERS
func Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]
	query_body := datasource.FindQuery(query)
	if query_body == nil {
		sendNotFound(w)
		return
	}
	sendJSON(w, query_body, http.StatusOK)
}

// Helper to set http.ResponseWriter
func setResponseContent(w http.ResponseWriter, r *http.Response, file string) http.ResponseWriter {
	if r.ContentLength > 1 {
		w.Header().Set("Content-Length", strconv.FormatInt(r.ContentLength, 10))
	}

	filename := file

	if reZip.MatchString(filename) {
		file_value := "attachment; filename=" + filename
		w.Header().Set("Content-Disposition", file_value)
	} else {
		file_value := "inline; filename=$" + filename
		w.Header().Set("Content-Disposition", file_value)
	}
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	return w
}

// HANDLERS TO TAKE CARE OF THE WEBPAGES
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Cat Game.\n")
}

// shows all purls
func PurlIndex(w http.ResponseWriter, r *http.Request) {
	ps := datasource.AllPurls()
	sendJSON(w, ps, http.StatusOK)
}

// admin interface
func AdminIndex(w http.ResponseWriter, r *http.Request) {
	ps := datasource.AllPurls()
	sendJSON(w, len(ps), http.StatusOK)
}

// gives back specific purl
func PurlShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purlId, err := strconv.Atoi(vars["purlId"])
	if err != nil {
		sendNotFound(w)
		return
	}
	purl := datasource.FindPurl(purlId)
	repoId, _ := strconv.Atoi(purl.Repo_obj_id)
	repo := datasource.FindRepoObj(repoId)
	if purl.Id == 0 || repo.Id == 0 {
		sendNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	M := struct {
		Id            int
		Information   string
		File_name     string
		Repo_url      string
		Repo_obj_id   string
		Last_accessed time.Time
		Access_count  int
	}{
		purl.Id,
		repo.Information,
		repo.Filename,
		repo.Url,
		purl.Repo_obj_id,
		purl.Last_accessed,
		purl.Access_count,
	}
	err = txViewTemplate.Execute(w, M)
	if err != nil {
		log.Println(err.Error())
	}
}

// Either copies over or redirects a file from a remote source
func PurlShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purlId, err := strconv.Atoi(vars["purlId"])
	if err != nil {
		sendNotFound(w)
		return
	}

	purl := datasource.FindPurl(purlId)
	if purl.Id == 0 {
		sendNotFound(w)
		return
	}

	repo_id, _ := strconv.Atoi(purl.Repo_obj_id)
	repo := datasource.FindRepoObj(repo_id)

	if repo.Id != repo_id {
		log.Println("Could not return correct repo object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if we cannot proxy the file redirect to it
	if regexpCurate.MatchString(repo.Information) {
		datasource.LogRecordAccess(r, repo.Id, purl.Id)
		http.Redirect(w, r, repo.Url, 302)
		return
	}

	// checks for fedora configuration information in env
	fedorausername := os.Getenv("FEDORA_USER")
	fedorapassword := os.Getenv("FEDORA_PASS")
	var back_end_new = repo.Url
	if fedorausername != "" && fedorapassword != "" {
		repl := `http$1://` + fedorausername + `:` + fedorapassword + `$2`
		back_end_new = reHttp.ReplaceAllString(repo.Url, repl)
	}
	resp, err := http.Get(back_end_new)
	if err != nil {
		log.Println("Unable to grab url:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		http.Error(w, "Content Unavailable", http.StatusInternalServerError)
		return
	}

	datasource.LogRecordAccess(r, repo.Id, purl.Id)

	w = setResponseContent(w, resp, vars["filename"])

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println(err)
	}
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/purl/create
*/
func PurlCreate(w http.ResponseWriter, r *http.Request) {
	var purl Purl
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err.Error())
		return
	}
	if err := json.Unmarshal(body, &purl); err != nil {
		sendJSON(w, err, 422) // unprocessable entity
		return
	}

	datasource.CreatePurl(purl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	// should a response body be returned?
}
