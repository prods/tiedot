// Miscellaneous function handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"runtime"
)

// Flush and close all data files and shutdown the entire program.
func Shutdown(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	HttpDB.Close()
	os.Exit(0)
}

// Copy this database into destination directory.
func Dump(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var dest string

	if IsNewAPIRoute(r) {
		var jsonDoc struct {
			Destination string `json:"destination"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&jsonDoc)
		if err != nil {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, fmt.Sprintf("'%v' is not valid JSON document.", jsonDoc), 400)
			return
		}
		dest = jsonDoc.Destination
		if dest == "" {
			// TODO: Wrap Error in Object (JSON)
			http.Error(w, "No destination was provided.", 400)
			return
		}
	} else {
		// TODO: Remove once Old API is discontinued
		if !Require(w, r, "dest", &dest) {
			return
		}
	}
	if err := HttpDB.Dump(dest); err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
}

// Return server memory statistics.
func MemStats(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	resp, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "Cannot serialize MemStats to JSON.", 500)
		return
	}
	w.Write(resp)
}

// Return server protocol version number.
func Version(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte("6"))
}
