package modules

import "github.com/julienschmidt/httprouter"

const (
	HTTP_METHOD_GET = iota
	HTTP_METHOD_POST = iota
	HTTP_METHOD_PUT = iota
	HTTP_METHOD_PATCH = iota
	HTTP_METHOD_DELETE = iota
)

type APIRoute struct {
	Path string
	Method int
	Handler httprouter.Handle
}

type IAPIModule interface {
	GetRoutes() []APIRoute
}

func RegisterAPIModule(router httprouter.Router, module IAPIModule) {
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
}