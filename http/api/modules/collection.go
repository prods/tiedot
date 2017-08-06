package modules


import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
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

func (module CollectionAPIModule) CreateNewCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Create New Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module CollectionAPIModule) RenameCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Rename Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module CollectionAPIModule) DropCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Drop Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module CollectionAPIModule) GetCollections(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Get All Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module CollectionAPIModule) SyncCollections(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Sync Collection Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}