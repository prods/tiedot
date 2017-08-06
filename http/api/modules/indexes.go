package modules

import (
	"net/http"
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	"errors"
	"strings"
)

type IndexesAPIModule struct {
	routes []APIRoute
	db *db.DB
}

type IndexPath struct {
	Path string `json:"path"`
}

func NewIndexesAPIModule(db *db.DB) *IndexesAPIModule {
	newInstance := new(IndexesAPIModule)
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/collection/:collection_name/index", newInstance.CreateNewIndex, true),
		DELETE("/collection/:collection_name/index", newInstance.RemoveIndex, true),
		GET("/collection/:collection_name/indexes", newInstance.GetIndexes, true),
	}
	return newInstance
}

func (index IndexesAPIModule) GetRoutes() []APIRoute {
	return index.routes
}

func (index IndexesAPIModule) GetName() string {
	return "Indexes"
}

// Put an index on a document path.
func (index IndexesAPIModule) CreateNewIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody IndexPath

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("create index", errors.New("No Collection was provided")))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "create index", func(doc interface{}) map[string]interface{} {
		if requestBody.Path == "" {
			return GetCollectionErrorObject("create index", errors.New("No Index path was provided."), collectionName)
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	// Perform Index Creation
	dbcol := index.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetIndexErrorObject("create index", errors.New("Collection does not exist"), collectionName, requestBody.Path))
		return
	}
	if err := dbcol.Index(strings.Split(requestBody.Path, ",")); err != nil {
		RespondWithBadRequest(w, GetIndexErrorObject("create index", err, collectionName, requestBody.Path))
		return
	}

	// Report Success on Completion
	RespondCreated(w, map[string]interface{} {
		"operation" : "create index",
		"collection": collectionName,
		"path" : requestBody.Path,
		"done" : true,
	})
	return
}

// Return all indexed paths.
func (index IndexesAPIModule) GetIndexes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get indexes", errors.New("No Collection was provided"), ""))
		return
	}

	// Perform Index Creation
	dbcol := index.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get indexes",errors.New("Collection does not exist"), collectionName))
		return
	}
	indexes := make([][]string, 0)
	for _, path := range dbcol.AllIndexes() {
		indexes = append(indexes, path)
	}
	resp, err := json.Marshal(indexes)
	if err != nil {
		RespondWithInternalError(w, GetCollectionErrorObject("get indexes", err, collectionName))
		return
	}
	w.Write(resp)
}

// Remove an indexed path.
func (index IndexesAPIModule) RemoveIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody IndexPath

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("remove index", errors.New("No Collection was provided")))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "remove index", func(doc interface{}) map[string]interface{} {
		if requestBody.Path == "" {
			return GetCollectionErrorObject("remove index", errors.New("No Index path was provided."), collectionName)
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	// Perform Index Creation
	dbcol := index.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetIndexErrorObject("remove index", errors.New("Collection does not exist"), collectionName, requestBody.Path))
		return
	}
	if err := dbcol.Unindex(strings.Split(requestBody.Path, ",")); err != nil {
		RespondWithBadRequest(w, GetIndexErrorObject("remove index", err, collectionName, requestBody.Path))
		return
	}

	// Report Success on Completion
	RespondCreated(w, map[string]interface{} {
		"operation" : "remove index",
		"collection": collectionName,
		"path" : requestBody.Path,
		"done" : true,
	})
	return
}