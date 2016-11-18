package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Cat Game.\n")
}

func PurlIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// b, _ := json.Marshal(datasource.AllPurls())
	if err := json.NewEncoder(w).Encode(datasource.AllPurls()); err != nil {
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
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// t := template.New("view.html")
		t, err := template.ParseFiles("./view.html")
		if err != nil {
			panic(err)
		}
		fmt.Println(t)
		err = t.Execute(w, purl)
		if err != nil {
			panic(err)
		}
		// if err := json.NewEncoder(w).Encode(purl); err != nil {
		// 	panic(err)
		// }
		return
	}
	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var query string
	query = vars["query"]
	query_body := datasource.FindQuery(query)
	if query_body != nil {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		t := template.New("search.html").Funcs(templateFuncs)
		t, _ = t.ParseFiles("search.html")
		err := t.Execute(w, query_body)
		if err != nil {
			panic(err)
		}
		// if err := json.NewEncoder(w).Encode(query_body); err != nil {
		// 	panic(err)
		// }
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func PurlShowFile(w http.ResponseWriter, r *http.Request) {
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

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

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

	t := datasource.CreatePurl(purl)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
// taken from http://stackoverflow.com/a/19995976

var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}
