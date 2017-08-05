package middlewares

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

// JWTAuth JWT Authorization Middleware
func JWTAuth(handler httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {



		handler(w, r, p)
	})
}
