package modules

import (
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"github.com/HouzuoGuo/tiedot/http/api/middlewares"
	"net/http"
	"encoding/json"
	"errors"
)

const (
	HTTP_METHOD_GET = iota
	HTTP_METHOD_POST = iota
	HTTP_METHOD_PUT = iota
	HTTP_METHOD_PATCH = iota
	HTTP_METHOD_DELETE = iota
)

const JSON_EOF_ERROR string = "EOF"

type APIRoute struct {
	Path string
	Method int
	Handler httprouter.Handle
}

type IAPIModule interface {
	GetRoutes() []APIRoute
	GetName() string
}

func Mount(router *httprouter.Router, module IAPIModule) {
	for _, r := range module.GetRoutes() {
		switch r.Method {
		case HTTP_METHOD_GET:
			router.GET(r.Path, r.Handler)
		case HTTP_METHOD_POST:
			router.POST(r.Path, r.Handler)
		case HTTP_METHOD_PUT:
			router.PUT(r.Path, r.Handler)
		case HTTP_METHOD_PATCH:
			router.PATCH(r.Path, r.Handler)
		case HTTP_METHOD_DELETE:
			router.DELETE(r.Path, r.Handler)
		}
	}
	tdlog.Noticef("+ API Module '%s' was mounted.", module.GetName())
}

func NewAPIRoute(method int, path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	routeHandler := middlewares.StandardResponse(handler)
	if requiresAuthentication {
		routeHandler = middlewares.JWTAuth(routeHandler)
	}

	newInstance := APIRoute{
		Method:method,
		Path:path,
		Handler: routeHandler,
	}
	return newInstance
}

func GET(path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	if requiresAuthentication {
		return NewAPIRoute(HTTP_METHOD_GET, path, handler, requiresAuthentication)
	} else {
		return NewAPIRoute(HTTP_METHOD_GET, path, handler, requiresAuthentication)
	}
}

func POST(path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	return NewAPIRoute(HTTP_METHOD_POST, path, handler, requiresAuthentication)
}

func PUT(path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	return NewAPIRoute(HTTP_METHOD_PUT, path, handler, requiresAuthentication)
}

func DELETE(path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	return NewAPIRoute(HTTP_METHOD_DELETE, path, handler, requiresAuthentication)
}

func PATCH(path string, handler httprouter.Handle, requiresAuthentication bool) APIRoute {
	return NewAPIRoute(HTTP_METHOD_PATCH, path, handler, requiresAuthentication)
}

func Respond(w http.ResponseWriter, code int, responsePayload map[string]interface{}) {
	m, _ := json.Marshal(responsePayload)
	w.WriteHeader(code)
	w.Write(m)
}

func RespondOk(w http.ResponseWriter, responsePayload map[string]interface{}) {
	m, _ := json.Marshal(responsePayload)
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func RespondCreated(w http.ResponseWriter, responsePayload map[string]interface{}) {
	m, _ := json.Marshal(responsePayload)
	w.WriteHeader(http.StatusCreated)
	w.Write(m)
}

func RespondWithInternalError(w http.ResponseWriter, responsePayload map[string]interface{}) {
	m, _ := json.Marshal(responsePayload)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(m)
}

func RespondWithBadRequest(w http.ResponseWriter, responsePayload map[string]interface{}) {
	m, _ := json.Marshal(responsePayload)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(m)
}

func GetEngineErrorObject(operation string, err error) map[string]interface{} {
	return map[string]interface{} {
		"operation" : operation,
		"error" : err.Error(),
	}
}

func GetCollectionErrorObject(operation string, err error, collectionName string) map[string]interface{} {
	return map[string]interface{} {
		"operation" : operation,
		"error" : err.Error(),
		"collection" : collectionName,
	}
}

func GetDocumentErrorObject(operation string, err error, collectionName string, documentId string) map[string]interface{} {
	return map[string]interface{} {
		"operation" : operation,
		"error" : err.Error(),
		"collection" : collectionName,
		"document": documentId,
	}
}

func GetIndexErrorObject(operation string, err error, collectionName string, indexPath string) map[string]interface{} {
	return map[string]interface{} {
		"operation" : operation,
		"error" : err.Error(),
		"collection" : collectionName,
		"index": indexPath,
	}
}

func GetCompletionObject(operation string) map[string]interface{} {
	return map[string]interface{} {
		"operation" : operation,
		"completed" : true,
	}
}

func ExtractAndValidateRequestBody(r *http.Request, payloadObject interface{}, operation string, validation func(interface{}) map[string]interface{}) map[string]interface{} {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payloadObject)
	if err != nil {
		if err.Error() == JSON_EOF_ERROR {
			return GetEngineErrorObject(operation, errors.New("No destination JSON object was provided."))
		}
		return GetEngineErrorObject(operation, err)
	}
	if validation != nil {
		return validation(&payloadObject)
	}
	return nil
}