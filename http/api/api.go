package api

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"net/http"
	"fmt"
	"github.com/HouzuoGuo/tiedot/http/api/middlewares"
)

var databseInstance *db.DB // HTTP API endpoints operate on this database

// Start Starts API Server
func Start(dir string, port int, tlsCrt, tlsKey, jwtPubKey, jwtPrivateKey, bind, authToken string) {
	// Database Instance Initialization
	var err error
	databseInstance, err = db.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	// Router Initialization
	router := httprouter.New()

	// New API
	// These endpoints are always available and do not require authentication
	router.GET("/", middlewares.StandardResponse(Welcome))
	router.GET("/version", middlewares.StandardResponse(Version))
	router.GET("/memstats", middlewares.StandardResponse(MemStats))

	// collection management (stop-the-world)
	router.POST("/collection/:collection_name", middlewares.JWTAuth(middlewares.StandardResponse(Create)))
	router.PUT("/collection/:collection_name/rename/:new_collection_name", middlewares.JWTAuth(middlewares.StandardResponse(Rename)))
	router.DELETE("/collection/:collection_name", middlewares.JWTAuth(middlewares.StandardResponse(Drop)))
	router.GET("/collections", middlewares.JWTAuth(middlewares.StandardResponse(All)))
	router.POST("/collection/:collection_name/scrub", middlewares.JWTAuth(middlewares.StandardResponse(Scrub)))
	router.POST("/sync", middlewares.JWTAuth(middlewares.StandardResponse(Sync)))
	// query
	router.POST("/collection/:collection_name/query", middlewares.JWTAuth(middlewares.StandardResponse(Query)))
	router.POST("/collection/:collection_name/count", middlewares.JWTAuth(middlewares.StandardResponse(Count)))
	// document management
	router.POST("/collection/:collection_name/doc",middlewares.JWTAuth(middlewares.StandardResponse(Insert)))
	router.GET("/collection/:collection_name/doc/:id", middlewares.JWTAuth(middlewares.StandardResponse(Get)))
	router.PUT("/collection/:collection_name/doc/:id", middlewares.JWTAuth(middlewares.StandardResponse(Update)))
	router.DELETE("/collection/:collection_name/doc/:id", middlewares.JWTAuth(middlewares.StandardResponse(Delete)))
	router.GET("/collection/:collection_name/page/:page/of/:total", middlewares.JWTAuth(middlewares.StandardResponse(GetPage)))
	// TODO: Review if it will make more sense for it to be just /collection/:collection_name/count
	router.GET("/collection/:collection_name/count/approx", middlewares.JWTAuth(middlewares.StandardResponse(ApproxDocCount)))
	// index management (stop-the-world)
	router.POST("/collection/:collection_name/index", middlewares.JWTAuth(middlewares.StandardResponse(Index)))
	router.DELETE("/collection/:collection_name/index", middlewares.JWTAuth(middlewares.StandardResponse(Unindex)))
	router.GET("/collection/:collection_name/indexes", middlewares.JWTAuth(middlewares.StandardResponse(Indexes)))
	// misc (stop-the-world)
	router.POST("/shutdown", middlewares.JWTAuth(middlewares.StandardResponse(Shutdown)))
	router.POST("/dump", middlewares.JWTAuth(middlewares.StandardResponse(Dump)))


	// Server Initialization

	iface := "all interfaces"
	if bind != "" {
		iface = bind
	}

	if tlsCrt != "" {
		tdlog.Noticef("Will listen on %s (HTTPS), port %d.", iface, port)
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", bind, port), tlsCrt, tlsKey, router); err != nil {
			tdlog.Panicf("Failed to start HTTPS service - %s", err)
		}
	} else {
		tdlog.Noticef("Will listen on %s (HTTP), port %d.", iface, port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), router)
	}

}
