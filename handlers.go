package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

var regexp_curate *regexp.Regexp = regexp.MustCompile(`^(CurateND - |Reformatting Unit:)`)
var re_http *regexp.Regexp = regexp.MustCompile(`http(s?)://(.+)`)
var re_zip *regexp.Regexp = regexp.MustCompile(`\b(ovf$)|\b(zip$)|\b(vmdk$)`)

// HELPERS FOR THE HANDLERS
func Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]
	query_body := datasource.FindQuery(query)
	if query_body != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(query_body); err != nil {
			log.Println(err.Error())
			return
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		log.Println(err.Error())
		return
	}
}

// Helper to set http.ResponseWriter
func setResponseContent(w http.ResponseWriter, r *http.Response, file string) http.ResponseWriter {
	if r.ContentLength > 1 {
		w.Header().Set("Content-Length", strconv.FormatInt(r.ContentLength, 10))
	} else if r.ContentLength < 0 {
		// fedora does not handle large file sizes correctly in HEAD requests to the legacy API.
		// un 2s-complement the number for it.
		con_length := (int64(1) << 32) + int64(r.ContentLength)
		w.Header().Set("Content-Length", strconv.FormatInt(con_length, 10))
	}

	filename := file

	if re_zip.MatchString(filename) {
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

func PurlIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	ps := datasource.AllPurls()
	if err := json.NewEncoder(w).Encode(ps); err != nil {
		log.Println(err.Error())
		return
	}
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	ps := datasource.AllPurls()
	if err := json.NewEncoder(w).Encode(len(ps)); err != nil {
		log.Println(err.Error())
		return
	}
}

func PurlShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var purlId int
	var err error
	if purlId, err = strconv.Atoi(vars["purlId"]); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			log.Println(err.Error())
			return
		}
	}
	purl := datasource.FindPurl(purlId)
	if purl.Id == 0 {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			log.Println(err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(purl); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// Either copies over or redirects a file from a remote source
func PurlShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		purlId int
		err    error
	)

	if purlId, err = strconv.Atoi(vars["purlId"]); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			log.Println(err.Error())
			return
		}
	}

	purl := datasource.FindPurl(purlId)

	if purl.Id == 0 {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			log.Println(err.Error())
			return
		}
	}

	repo_id, _ := strconv.Atoi(purl.Repo_obj_id)
	repo := datasource.FindRepoObj(repo_id)

	if repo.Id != repo_id {
		log.Println("Could not return correct repo object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if we cannot proxy the file redirect to it
	if regexp_curate.MatchString(repo.Information) {
		datasource.LogRecordAccess(r, repo.Id, purl.Id)
		http.Redirect(w, r, repo.Url, 302)
		return
	}

	// checks for fedora configuration information in env
	fedorausername := os.Getenv("FEDORA_USER")
	fedorapassword := os.Getenv("FEDORA_PASS")
	var back_end_new string
	if fedorausername != "" && fedorapassword != "" {
		repl := `http$1://` + fedorausername + `:` + fedorapassword + `$2`
		back_end_new = re_http.ReplaceAllString(repo.Url, repl)
	} else {
		repl := `http$1://$2`
		back_end_new = re_http.ReplaceAllString(repo.Url, repl)
	}
	resp, err := http.Get(back_end_new)
	log.Println(back_end_new)
	if err != nil {
		log.Println("Unable to grab file: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	datasource.LogRecordAccess(r, repo.Id, purl.Id)

	w = setResponseContent(w, resp, vars["filename"])

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println(err.Error())
			return
		}
	}

	datasource.CreatePurl(purl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}
