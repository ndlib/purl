package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

// HELPERS FOR THE HANDLERSX
func Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var query string
	query = vars["query"]
	query_body := datasource.FindQuery(query)
	if query_body != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(query_body); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

// Helper to set http.ResponseWriter
func setResponseContent(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	vars := mux.Vars(r)

	if r.ContentLength > 1 {
		w.Header().Set("Content-Length", strconv.FormatInt(r.ContentLength, 10))
	} else if r.ContentLength < 0 {
		con_length := int64(math.Pow(float64(2), float64(32))) + r.ContentLength
		w.Header().Set("Content-Length", strconv.FormatInt(con_length, 10))
	}

	filename := vars["filename"]
	re := regexp.MustCompile(`\b(ovf$)|\b(zip$)|\b(vmdk$)`)
	if re.MatchString(filename) {
		file_value := "attachment; filename=$" + filename
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
		panic(err)
	}
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	ps := datasource.AllPurls()
	if err := json.NewEncoder(w).Encode(len(ps)); err != nil {
		panic(err)
	}
}

func PurlShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var purlId int
	var err error
	if purlId, err = strconv.Atoi(vars["purlId"]); err != nil {
		panic(err)
	}
	purl := datasource.FindPurl(purlId)
	if purl.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(purl); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

// Either copies over or redirects a file from a remote source
func PurlShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		purlId int
		err    error
	)

	if purlId, err = strconv.Atoi(vars["purlId"]); err != nil {
		log.Println("purlID err:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purl := datasource.FindPurl(purlId)

	if purl.Id > 0 {
		repo_id, _ := strconv.Atoi(purl.Repo_obj_id)
		repo := datasource.FindRepoObj(repo_id)

		if repo.Id != repo_id {
			log.Printf("Could not return correct repo object")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if we cannot proxy the file redirect to it
		if regexp.MustCompile(`^(CurateND - |Reformatting Unit:)`).MatchString(repo.Information) {
			datasource.LogRecordAccess(r, repo.Id, purl.Id)
			http.Redirect(w, r, repo.Url, 302)
			return
		}

		// TODO: Add permissions to our url request
		re := regexp.MustCompile(`http(s):\/\/(.+)`)
		// fedorausername := "user"
		// fedorapassword := "pass"
		repl := `http$1://$2` //` + fedorausername + `:` + fedorapassword + `
		back_end_new := re.ReplaceAllString(repo.Url, repl)

		resp, err := http.Get(back_end_new)
		if err != nil {
			log.Printf("Unable to grab file: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		datasource.LogRecordAccess(r, repo.Id, purl.Id)

		w = setResponseContent(w, r)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
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
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &purl); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	datasource.CreatePurl(purl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}
