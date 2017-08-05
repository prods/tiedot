package modules

import (
	"net/http"
	"os"
	"fmt"
	"runtime"
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
)

type MiscellaneousAPIModule struct {
	routes []APIRoute
	db *db.DB
}

func NewMiscellaneousAPIModule(db *db.DB) *MiscellaneousAPIModule {
	newInstance := new(MiscellaneousAPIModule)
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/shutdown", newInstance.Shutdown, true),
		POST("/dump", newInstance.Dump, true),
		GET("/memstats", newInstance.MemStats, false),
		GET("/version",newInstance.Version, false),
	}
	return newInstance
}

func (misc MiscellaneousAPIModule) GetRoutes() []APIRoute {
	return misc.routes
}

func (misc MiscellaneousAPIModule) GetName() string {
	return "Miscellaneous"
}

// Flush and close all data files and shutdown the entire program.
func (misc MiscellaneousAPIModule) Shutdown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	misc.db.Close()
	os.Exit(0)
}

// Copy this database into destination directory.
func (misc MiscellaneousAPIModule) Dump(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var dest string
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
	if err := misc.db.Dump(dest); err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
}

// Return server memory statistics.
func (misc MiscellaneousAPIModule) MemStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
func (misc MiscellaneousAPIModule) Version(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("6"))
}