package modules


import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
	"errors"
	"encoding/json"
)

type CollectionAPIModule struct {
	*APIModuleBase
}

func NewCollectionAPIModule(db *db.DB) *CollectionAPIModule {
	newInstance := &CollectionAPIModule{&APIModuleBase{}}
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/collection/:collection_name", newInstance.CreateNewCollection, true),
		PUT("/collection/:collection_name", newInstance.RenameCollection, true),
		DELETE("/collection/:collection_name", newInstance.DropCollection, true),
		POST("/collection/:collection_name/scrub", newInstance.ScrubCollection, true),
		GET("/collections", newInstance.GetCollections, true),
		GET("/sync", newInstance.SyncCollections, true),
	}
	return newInstance
}

func (module CollectionAPIModule) GetRoutes() []APIRoute {
	return module.routes
}

func (module CollectionAPIModule) GetName() string {
	return "Collections"
}

func (module CollectionAPIModule) GetDocumentation() APIModuleDocumentation {
	return module.documentation
}

// CreateNewCollection Creates a new Collection
func (module CollectionAPIModule) CreateNewCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("create index", errors.New("No Collection was provided")))
		return
	}
	if err := module.db.Create(collectionName); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("create collection", err, collectionName))
		return
	} else {
		RespondCreated(w, map[string]interface{} {
			"operation" : "create collection",
			"collection": collectionName,
			"done": true,
		})
	}
}

// RenameCollection Renames a specified Collection
func (module CollectionAPIModule) RenameCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody struct {
		NewName string 	`json:"new_name"`
	}
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("rename collection", errors.New("No Collection was provided")))
		return
	}
	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "rename collection", func(doc interface{}) map[string]interface{} {
		if requestBody.NewName == "" {
			return GetCollectionErrorObject("rename collection", errors.New("Collection new name was not provided."), collectionName)
		}
		return nil
	}); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	// Perform Collection rename
	if err := module.db.Rename(collectionName, requestBody.NewName); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("rename collection", err, collectionName))
		return
	}

	RespondOk(w, map[string]interface{} {
		"operation" : "rename collection",
		"collection": collectionName,
		"done": true,
	})
}

// DropCollection Drops a specified Collection
func (module CollectionAPIModule) DropCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("rename collection", errors.New("No Collection was provided")))
		return
	}
	// Perform Collection Drop
	if err := module.db.Drop(collectionName); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("drop collection", err, collectionName))
		return
	}

	RespondOk(w, map[string]interface{} {
		"operation" : "drop collection",
		"collection": collectionName,
		"done": true,
	})
}

// SCrubCollection Scrubs an specified collection content
func (module CollectionAPIModule) ScrubCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetEngineErrorObject("scrub collection", errors.New("No Collection was provided")))
		return
	}
	// Perform Collection Scrub
	dbCollection := module.db.Use(collectionName)
	if dbCollection == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("scrub collection", errors.New("Collection does not exist"), collectionName))
		return
	}

	if err := module.db.Scrub(collectionName); err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("scrub collection", err, collectionName))
	}

	RespondOk(w, map[string]interface{} {
		"operation" : "scrub collection",
		"collection": collectionName,
		"done": true,
	})
}

// GetCollections Gets all Collections
func (module CollectionAPIModule) GetCollections(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cols := make([]string, 0)
	for _, v := range module.db.AllCols() {
		cols = append(cols, v)
	}
	resp, err := json.Marshal(cols)
	if err != nil {
		RespondWithInternalError(w, GetEngineErrorObject("get all collections", err))
		return
	}
	w.Write(resp)
}

// SyncCollection Syncs an specified collection
func (module CollectionAPIModule) SyncCollections(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Sync Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}