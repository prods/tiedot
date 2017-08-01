package httpapi

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

func StandardAPIResponder(handler httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Cache-Control", "must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		handler(w, r, p)
	})
}

func Respond(w http.ResponseWriter, message map[string]interface{}, code int) {
	m, _ := json.Marshal(message)
	w.WriteHeader(code)
	w.Write(m)
}

func RespondWithEngineError(w http.ResponseWriter, operation string, message string, err error) {
	Respond(w, map[string]interface{}{
		"operation" : operation,
		"error" : err.Error(),
	}, http.StatusInternalServerError)
}

func RespondWithCollectionError(w http.ResponseWriter, operation string, message string, err error, collectionName string) {
	Respond(w, map[string]interface{}{
		"operation" : operation,
		"error" : err.Error(),
		"collection" : collectionName,
	}, http.StatusInternalServerError)
}

func RespondWithDocumentError(w http.ResponseWriter, operation string, message string, err error, collectionName string, documentId string) {
	Respond(w, map[string]interface{}{
		"operation" : operation,
		"error" : err.Error(),
		"collection" : collectionName,
		"document": documentId,
	}, http.StatusInternalServerError)
}
