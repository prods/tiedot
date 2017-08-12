package modules

import (
	"github.com/HouzuoGuo/tiedot/db"

	"net/http"
	"github.com/julienschmidt/httprouter"
	"errors"
	"strconv"
	"encoding/json"
	"fmt"
)

type DocumentAPIModule struct {
	*APIModuleBase
}

func NewDocumentAPIModule(db *db.DB) *DocumentAPIModule {
	newInstance := &DocumentAPIModule{&APIModuleBase{}}
	newInstance.db = db
	newInstance.routes = []APIRoute{
		POST("/collection/:collection_name/doc", newInstance.InsertDocument, true),
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

func (module DocumentAPIModule) InsertDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody map[string]interface{}

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("insert document", errors.New("No Collection was provided"), ""))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "insert document", nil); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("insert document", errors.New("Collection does not exist"), collectionName))
		return
	}

	id, err := dbcol.Insert(requestBody)
	if err != nil {
		RespondWithInternalError(w, GetCollectionErrorObject("insert document", err, collectionName))
		return
	}

	RespondCreated(w, map[string]interface{}{
		"operation":  "create document",
		"collection": collectionName,
		"document":   id,
		"done":       true,
	})
	return
}

func (module DocumentAPIModule) UpdateDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestBody map[string]interface{}

	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("update document", errors.New("No Collection was provided"), ""))
		return
	}

	// Get Document Id
	documentIdValue := p.ByName("id")
	if documentIdValue == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("update document", errors.New("No Document Id was provided"), collectionName))
		return
	}

	documentId, err := strconv.Atoi(documentIdValue)
	if err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("update document", errors.New("Invalid Document Id"), collectionName))
		return
	}

	// Extract and Validate Request Body
	if apiErr := ExtractAndValidateRequestBody(r, &requestBody, "update document", nil); apiErr != nil {
		RespondWithBadRequest(w, apiErr)
		return
	}

	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("update document", errors.New("Collection does not exist"), collectionName))
		return
	}

	if err := dbcol.Update(documentId, requestBody); err != nil {
		RespondWithInternalError(w, GetDocumentErrorObject("update document", err, collectionName, documentId))
		return
	}

	RespondCreated(w, map[string]interface{}{
		"operation":  "update document",
		"collection": collectionName,
		"document":   documentId,
		"done":       true,
	})
	return
}

func (module DocumentAPIModule) DeleteDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("delete document", errors.New("No Collection was provided"), ""))
		return
	}

	// Get Document Id
	documentIdValue := p.ByName("id")
	if documentIdValue == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("delete document", errors.New("No Document Id was provided"), collectionName))
		return
	}

	documentId, err := strconv.Atoi(documentIdValue)
	if err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("delete document", errors.New("Invalid Document Id"), collectionName))
		return
	}

	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("delete document", errors.New("Collection does not exist"), collectionName))
		return
	}

	if err := dbcol.Delete(documentId); err != nil {
		RespondWithInternalError(w, GetDocumentErrorObject("delete document", err, collectionName, documentId))
		return
	}

	RespondCreated(w, map[string]interface{}{
		"operation":  "delete document",
		"collection": collectionName,
		"document":   documentId,
		"done":       true,
	})
	return
}

func (module DocumentAPIModule) GetDocument(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get document", errors.New("No Collection was provided"), ""))
		return
	}

	// Get Document Id
	documentIdValue := p.ByName("id")
	if documentIdValue == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get document", errors.New("No Document Id was provided"), collectionName))
		return
	}

	// Get Document Id Integer Value
	documentId, err := strconv.Atoi(documentIdValue)
	if err != nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get document", errors.New("Invalid Document Id"), collectionName))
		return
	}

	// Use Collection
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get document", errors.New("Collection does not exist"), collectionName))
		return
	}

	// Get Document from Collection
	document, err := dbcol.Read(documentId)
	if err != nil {
		RespondWithInternalError(w, GetDocumentErrorObject("get document", err, collectionName, documentId))
		return
	}
	if document == nil {
		RespondNotFound(w, GetDocumentErrorObject("get document", errors.New("No such document ID"), collectionName, documentId))
		return
	}

	// Return Document
	RespondOk(w, document)
}

func (module DocumentAPIModule) GetAllDocuments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("No Collection was provided"), ""))
		return
	}

	// Use Collection
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("Collection does not exist"), collectionName))
		return
	}

	// Get all Documents
	documents := make(map[string]interface{})
	dbcol.ForEachDoc(func(id int, doc []byte) bool {
		var docObj map[string]interface{}
		if err := json.Unmarshal(doc, &docObj); err == nil {
			documents[strconv.Itoa(id)] = docObj
		}
		return true
	})

	// Return Results
	RespondOk(w, documents)
}

func (module DocumentAPIModule) GetPageOfDocuments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get page of documents", errors.New("No Collection was provided"), ""))
		return
	}

	// Get Total per page value
	totalValue := p.ByName("total")
	if totalValue == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get page of documents", errors.New("No Total per page was provided"), ""))
		return
	}

	// Get Document Id Integer Value
	totalPerPage, err := strconv.Atoi(totalValue)
	if err != nil || totalPerPage < 1 {
		RespondWithBadRequest(w, GetCollectionErrorObject("get document", errors.New(fmt.Sprintf("Invalid Total per page %s", totalPerPage)), collectionName))
		return
	}

	// Use Collection
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get page of documents", errors.New("Collection does not exist"), collectionName))
		return
	}

	// Get Page Number Value
	pageValue := p.ByName("page")
	if pageValue == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get page of documents", errors.New("No Page Number was provided"), collectionName))
		return
	}

	// Get Page Number Integer
	pageNumber, err := strconv.Atoi(pageValue)
	if err != nil || pageNumber < 0 || pageNumber >= totalPerPage {
		RespondWithBadRequest(w, GetCollectionErrorObject("get page of documents", errors.New(fmt.Sprintf("Invalid Page Number %s", pageValue)), collectionName))
		return
	}

	// Extract Values from collection
	documents := make(map[string]interface{})
	dbcol.ForEachDocInPage(pageNumber, totalPerPage, func(id int, doc []byte) bool {
		var docObj map[string]interface{}
		if err := json.Unmarshal(doc, &docObj); err == nil {
			documents[strconv.Itoa(id)] = docObj
		}
		return true
	})

	// Return Results
	RespondOk(w, documents)
}

func (module DocumentAPIModule) GetApproximateCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get Collection Name
	collectionName := p.ByName("collection_name")
	if collectionName == "" {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("No Collection was provided"), ""))
		return
	}

	// Use Collection
	dbcol := module.db.Use(collectionName)
	if dbcol == nil {
		RespondWithBadRequest(w, GetCollectionErrorObject("get all documents", errors.New("Collection does not exist"), collectionName))
		return
	}

	RespondOk(w, map[string]interface{} {
		"collection": collectionName,
		"count":  dbcol.ApproxDocCount(),
	})
}
