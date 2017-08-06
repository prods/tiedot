package modules

import (
	"net/http"
	"os"
	"runtime"
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	"errors"
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
	RespondOk(w, GetCompletionObject("shutdown"))
	os.Exit(0)
}

// Copy this database into destination directory.
func (misc MiscellaneousAPIModule) Dump(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestBody struct {
		Destination string `json:"destination"`
	}
	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "dump", func(doc interface{}) map[string]interface{} {
		if requestBody.Destination == "" {
			return GetEngineErrorObject("dump", errors.New("No destination was provided."))
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}
	// Perform Operation
	if err := misc.db.Dump(requestBody.Destination); err != nil {
		RespondWithInternalError(w, GetEngineErrorObject("dump", err))
		return
	}
}

// Return server memory statistics.
func (misc MiscellaneousAPIModule) MemStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	resp, err := json.Marshal(stats)
	if err != nil {
		RespondWithInternalError(w, GetEngineErrorObject("memstat", errors.New("Cannot serialize MemStats to JSON.")))
		return
	}
	w.Write(resp)
}

// Return server protocol version number.
func (misc MiscellaneousAPIModule) Version(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	RespondOk(w, map[string]interface{} {
		"version" : map[string]interface{} {
			"engine" : "0",
			"api": 2,
			"protocol": 6,
		},
	})
}