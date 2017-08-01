// Index management handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type IndexPath struct {
	Path string `json:"path"`
}

// Put an index on a document path.
func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col, path string

	if IsNewAPIRoute(r) {
		var jsonDoc IndexPath
		col = p.ByName("collection_name")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&jsonDoc)
		if err != nil {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, fmt.Sprintf("'%v' is not valid JSON document.", jsonDoc), 400)
			return
		}
		path = jsonDoc.Path
		if path == "" {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, "No Index path was provided.", 400)
			return
		}
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
		if !Require(w, r, "path", &path) {
			return
		}
	}
	dbcol := HttpDB.Use(col)
	if dbcol == nil {
		http.Error(w, fmt.Sprintf("Collection '%s' does not exist.", col), 400)
		return
	}
	if err := dbcol.Index(strings.Split(path, ",")); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.WriteHeader(201)
}

// Return all indexed paths.
func Indexes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col string
	if IsNewAPIRoute(r) {
		col = p.ByName("collection_name")
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
	}
	dbcol := HttpDB.Use(col)
	if dbcol == nil {
		http.Error(w, fmt.Sprintf("Collection '%s' does not exist.", col), 400)
		return
	}
	indexes := make([][]string, 0)
	for _, path := range dbcol.AllIndexes() {
		indexes = append(indexes, path)
	}
	resp, err := json.Marshal(indexes)
	if err != nil {
		http.Error(w, fmt.Sprint("Server error."), 500)
		return
	}
	w.Write(resp)
}

// Remove an indexed path.
func Unindex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col, path string

	if IsNewAPIRoute(r) {
		var jsonDoc IndexPath
		col = p.ByName("collection_name")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&jsonDoc)
		if err != nil {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, fmt.Sprintf("'%v' is not valid JSON document.", jsonDoc), 400)
			return
		}
		path = jsonDoc.Path
		if path == "" {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, "No Index path was provided.", 400)
			return
		}
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
		if !Require(w, r, "path", &path) {
			return
		}
	}
	dbcol := HttpDB.Use(col)
	if dbcol == nil {
		http.Error(w, fmt.Sprintf("Collection '%s' does not exist.", col), 400)
		return
	}
	if err := dbcol.Unindex(strings.Split(path, ",")); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
}
