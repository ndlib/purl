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
	"strings"
	"time"

	"github.com/gorilla/mux"
)

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

var attachmentExt = regexp.MustCompile(`\b(ovf$)|\b(zip$)|\b(vmdk$)`)

// Helper to set http.ResponseWriter
func setResponseContent(w http.ResponseWriter, r *http.Response, filename string) {
	if r.ContentLength > 1 {
		w.Header().Set("Content-Length", strconv.FormatInt(r.ContentLength, 10))
	}

	// For certain file extensions, we set the download to be an "attachemnt"
	// so a web browser will not try to open it in the browser window.
	//
	// We use the filename passed in through the URL, and not the filename
	// stored in the purl record.
	//
	// All of this is previous behavior. Maybe it could be rethought out.
	var disposition string
	if attachmentExt.MatchString(filename) {
		disposition = "attachment; filename=" + filename
	} else {
		disposition = "inline; filename=" + filename
	}
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
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

var (
	redirectPattern = regexp.MustCompile(`^(CurateND - |Reformatting Unit:)`)
)

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

	// Some entries need a redirect and not a proxy. Determining that from
	// special patterns in the information string is legacy behavior.
	if redirectPattern.MatchString(repo.Information) {
		datasource.LogRecordAccess(r, repo.Id, purl.Id)
		http.Redirect(w, r, repo.Url, 302)
		return
	}

	proxyRequest, _ := http.NewRequest("GET", repo.Url, nil)
	// this test is a little hokey. we don't need to send auth to non-fedora
	// urls. We assume every fedora URL has the word "fedora" in it somewhere.
	if strings.Contains(repo.Url, "fedora") {
		// checks for fedora configuration information in env
		fedorausername := os.Getenv("FEDORA_USER")
		fedorapassword := os.Getenv("FEDORA_PASS")
		if fedorausername != "" && fedorapassword != "" {
			proxyRequest.SetBasicAuth(fedorausername, fedorapassword)
		}
	}
	resp, err := http.DefaultClient.Do(proxyRequest)
	if err != nil {
		log.Println("Unable to grab url:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("upstream returned", resp.StatusCode, repo.Url)
		http.Error(w, "Content Unavailable", http.StatusInternalServerError)
		return
	}

	datasource.LogRecordAccess(r, repo.Id, purl.Id)

	// this uses the filename that the client passed us...should we use
	// the filename stored in the purl record instead?
	setResponseContent(w, resp, vars["filename"])

	if r.ContentLength > 0 {
		_, err = io.CopyN(w, resp.Body, r.ContentLength)
	} else {
		_, err = io.Copy(w, resp.Body)
	}
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
