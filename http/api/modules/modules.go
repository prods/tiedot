package modules

import (
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"github.com/HouzuoGuo/tiedot/http/api/middlewares"
)

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
	tdlog.Noticef("API Module '%' was mounted.", module.GetName())
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