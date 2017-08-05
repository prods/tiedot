package modules

import (
	"net/http"
	"os"
	"fmt"
	"runtime"
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
)

type MiscellaneousAPIModule struct {
	routes []APIRoute
	db *db.DB
}

func NewMiscellaneousAPIModule(db *db.DB) *MiscellaneousAPIModule {
	newInstance := new(MiscellaneousAPIModule)
	newInstance.routes = make([]APIRoute, 0)
	return newInstance
}

func (misc MiscellaneousAPIModule) GetRoutes() []APIRoute {
	return misc.routes
}


// Flush and close all data files and shutdown the entire program.
func (misc MiscellaneousAPIModule) Shutdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	db.Close()
	os.Exit(0)
}

// Copy this database into destination directory.
func (misc MiscellaneousAPIModule) Dump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	var dest string
	if !Require(w, r, "dest", &dest) {
		return
	}
	if err := db.Dump(dest); err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
}

// Return server memory statistics.
func (misc MiscellaneousAPIModule) MemStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
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
func (misc MiscellaneousAPIModule) Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, PUT, OPTIONS")
	w.Write([]byte("6"))
}