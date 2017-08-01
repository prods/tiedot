// Collection management handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Create a collection.
func Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col string

	if IsNewAPIRoute(r) {
		col = p.ByName("collection_name")
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
	}
	if err := HttpDB.Create(col); err != nil {
		// TODO: Wrap Error in Object (JSON)
		http.Error(w, fmt.Sprint(err), 400)
	} else {
		w.WriteHeader(201)
	}
}

// Return all collection names.
func All(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cols := make([]string, 0)
	for _, v := range HttpDB.AllCols() {
		cols = append(cols, v)
	}
	resp, err := json.Marshal(cols)
	if err != nil {
		// TODO: Wrap Error in Object (JSON)
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	w.Write(resp)
}

// Rename a collection.
func Rename(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var oldName, newName string

	if IsNewAPIRoute(r) {
		oldName = p.ByName("collection_name")
		newName = p.ByName("new_collection_name")
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "old", &oldName) {
			return
		}
		if !Require(w, r, "new", &newName) {
			return
		}
	}

	if err := HttpDB.Rename(oldName, newName); err != nil {
		// TODO: Wrap Error in Object (JSON)
		http.Error(w, fmt.Sprint(err), 400)
	}
}

// Drop a collection.
func Drop(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col string
	if IsNewAPIRoute(r) {
		col = p.ByName("collection_name")
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
	}
	if err := HttpDB.Drop(col); err != nil {
		// TODO: Wrap Error in Object (JSON)
		http.Error(w, fmt.Sprint(err), 400)
	}
}

// De-fragment collection free space and fix corrupted documents.
func Scrub(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var col string
	if IsNewAPIRoute(r) {
		col = p.ByName("collection_name")
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "col", &col) {
			return
		}
	}
	dbCol := HttpDB.Use(col)
	if dbCol == nil {
		// TODO: Wrap Error in Object (JSON)
		http.Error(w, fmt.Sprintf("Collection %s does not exist", col), 400)
	} else {
		HttpDB.Scrub(col)
	}
}

/*
Noop
*/
func Sync(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
}
