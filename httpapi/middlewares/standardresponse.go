package middlewares

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// StandardResponse Standard Response Middleware. Ensure all common headers are setup
func StandardResponse(handler httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Cache-Control", "must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		handler(w, r, p)
	})
}