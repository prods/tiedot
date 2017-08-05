// Collection management handlers.

package legacy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/httpapi/shared"
)

// Create a collection.
func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
	if err := shared.GetDatabaseInstance().Create(col); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	} else {
		w.WriteHeader(201)
	}
}

// Return all collection names.
func All(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	cols := make([]string, 0)
	for _, v := range shared.GetDatabaseInstance().AllCols() {
		cols = append(cols, v)
	}
	resp, err := json.Marshal(cols)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	w.Write(resp)
}

// Rename a collection.
func Rename(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	var oldName, newName string
	if !Require(w, r, "old", &oldName) {
		return
	}
	if !Require(w, r, "new", &newName) {
		return
	}
	if err := shared.GetDatabaseInstance().Rename(oldName, newName); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	}
}

// Drop a collection.
func Drop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
	if err := shared.GetDatabaseInstance().Drop(col); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
	}
}

// De-fragment collection free space and fix corrupted documents.
func Scrub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
	dbCol := shared.GetDatabaseInstance().Use(col)
	if dbCol == nil {
		http.Error(w, fmt.Sprintf("Collection %s does not exist", col), 400)
	} else {
		shared.GetDatabaseInstance().Scrub(col)
	}
}

/*
Noop
*/
func Sync(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
}
