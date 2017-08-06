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
	*APIModuleBase
}

type IndexPath struct {
	Path []string `json:"path"`
}

func NewIndexesAPIModule(db *db.DB) *IndexesAPIModule {
	newInstance := &IndexesAPIModule{&APIModuleBase{}}
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/collection/:collection_name/index", newInstance.CreateNewIndex, true),
		DELETE("/collection/:collection_name/index", newInstance.RemoveIndex, true),
		GET("/collection/:collection_name/indexes", newInstance.GetIndexes, true),
	}
	return newInstance
}

func (module IndexesAPIModule) GetRoutes() []APIRoute {
	return module.routes
}

func (module IndexesAPIModule) GetName() string {
	return "Indexes"
}

func (module IndexesAPIModule) GetDocumentation() APIModuleDocumentation {
	return module.documentation
}

// Put an index on a document path.
func (module IndexesAPIModule) CreateNewIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody IndexPath

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("create index", errors.New("No Collection was provided")))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "create index", func(doc interface{}) map[string]interface{} {
		if len(requestBody.Path) == 0 {
			return GetCollectionErrorObject("create index", errors.New("No Index path was provided."), collectionName)
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	// Perform Index Creation
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetIndexErrorObject("create index", errors.New("Collection does not exist"), collectionName, strings.Join(requestBody.Path, ",")))
		return
	}
	if err := dbcol.Index(requestBody.Path); err != nil {
		RespondWithBadRequest(w, GetIndexErrorObject("create index", err, collectionName, strings.Join(requestBody.Path, ",")))
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
func (module IndexesAPIModule) GetIndexes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get indexes", errors.New("No Collection was provided"), ""))
		return
	}

	// Perform Index Creation
	dbcol := module.db.Use(collectionName)
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
func (module IndexesAPIModule) RemoveIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody IndexPath

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("remove index", errors.New("No Collection was provided")))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "remove index", func(doc interface{}) map[string]interface{} {
		if len(requestBody.Path) == 0 {
			return GetCollectionErrorObject("remove index", errors.New("No Index path was provided."), collectionName)
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	// Perform Index Creation
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetIndexErrorObject("remove index", errors.New("Collection does not exist"), collectionName, strings.Join(requestBody.Path, ",")))
		return
	}
	if err := dbcol.Unindex(requestBody.Path); err != nil {
		RespondWithBadRequest(w, GetIndexErrorObject("remove index", err, collectionName, strings.Join(requestBody.Path, ",")))
		return
	}

	// Report Success on Completion
	RespondOk(w, map[string]interface{} {
		"operation" : "remove index",
		"collection": collectionName,
		"path" : requestBody.Path,
		"done" : true,
	})
	return
}