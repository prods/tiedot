package modules


import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
)

type DocumentAPIModule struct {
	*APIModuleBase
}

func NewDocumentAPIModule(db *db.DB) *DocumentAPIModule {
	newInstance := &DocumentAPIModule{&APIModuleBase{}}
	newInstance.db = db
	newInstance.routes = []APIRoute {
		POST("/collection/:collection_name/doc", newInstance.CreateNewDocument, true),
		PUT("/collection/:collection_name/doc/:id", newInstance.UpdateDocument, true),
		DELETE("/collection/:collection_name/doc/:id", newInstance.DeleteDocument, true),
		GET("/collection/:collection_name/doc/:id", newInstance.GetDocument, true),
		GET("/collection/:collection_name/page/:page/of/:total", newInstance.GetPageOfDocuments, true),
		GET("/collection/:collection_name/docs", newInstance.GetAllDocuments, true),
		GET("/collection/:collection_name/docs/count", newInstance.GetApproximateCount, true),
	}
	return newInstance
}

func (module DocumentAPIModule) GetRoutes() []APIRoute {
	return module.routes
}

func (module DocumentAPIModule) GetName() string {
	return "Documents"
}

func (module DocumentAPIModule) GetDocumentation() APIModuleDocumentation {
	return module.documentation
}

func (module DocumentAPIModule) CreateNewDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Create New Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) UpdateDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Update Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) DeleteDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Remove Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) GetDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Get Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) GetAllDocuments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Get All Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) GetPageOfDocuments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Get Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}

func (module DocumentAPIModule) GetApproximateCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: Pending Get Document Implementation
	Respond(w, http.StatusNotImplemented, map[string]interface{} {
		"status" : "This operation is not yet supported on this endpoint.",
	})
}