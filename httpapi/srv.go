/*
Register HTTP API endpoints and handle authorization requirements.

Without specifying authorization parameters in the command line, tiedot server does not
require any authorization on any endpoint.

tiedot supports two authorization mechanisms:
- Pre-shared authorization token
The API endpoints will require 'Authorization: token PRE_SHARED_TOKEN' header. The pre-shared
token is specified in command line parameter "-authtoken".
Client request example: curl -I -H "Authorization: token PRE_SHARED_TOKEN" http://127.0.0.1:8080/all
- JWT (JSON Web Token)
The sophisticated mechanism offers finer-grained access control, separated by individual users.
Access to specific endpoints are granted explicitly to each user.

These API endpoints will never require authorization: / (root), /version, and /memstats
*/

package httpapi

import (
	"fmt"
	"net/http"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"github.com/HouzuoGuo/tiedot/httpapi/legacy"
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/httpapi/middlewares"
	"github.com/HouzuoGuo/tiedot/httpapi/shared"
)

// Start HTTP server and block until the server shuts down. Panic on error.
func Start(dir string, port int, tlsCrt, tlsKey, jwtPubKey, jwtPrivateKey, bind, authToken string, supportLegacyAPI bool) {
	// Initialize Database Instance
	err := shared.InitializeDatabaseInstance(dir)
	if err != nil {
		panic(err)
	}
	
	// Initialize JWT Support 
	jwtSupport, err := middlewares.NewJWTAuthentication(jwtPubKey, jwtPrivateKey, authToken)
	if err != nil {
		panic(err)
	}
	jwtAuthMiddleware := jwtSupport.AuthenticationHandler
	
	// Initialize Router
	router := httprouter.New()

	if supportLegacyAPI {
		// These endpoints are always available and do not require authentication
		router.GET("/", Welcome)
		router.GET("/version", legacy.Version)
		router.GET("/memstats", legacy.MemStats)
		// collection management (stop-the-world)
		router.GET("/create", jwtAuthMiddleware(legacy.Create))
		router.GET("/rename", jwtAuthMiddleware(legacy.Rename))
		router.GET("/drop", jwtAuthMiddleware(legacy.Drop))
		router.GET("/all", jwtAuthMiddleware(legacy.All))
		router.GET("/scrub", jwtAuthMiddleware(legacy.Scrub))
		router.GET("/sync", jwtAuthMiddleware(legacy.Sync))
		// query
		router.GET("/query", jwtAuthMiddleware(legacy.Query))
		router.GET("/count", jwtAuthMiddleware(legacy.Count))
		// document management
		router.GET("/insert", jwtAuthMiddleware(legacy.Insert))
		router.GET("/get", jwtAuthMiddleware(legacy.Get))
		router.GET("/getpage", jwtAuthMiddleware(legacy.GetPage))
		router.GET("/update", jwtAuthMiddleware(legacy.Update))
		router.GET("/delete", jwtAuthMiddleware(legacy.Delete))
		router.GET("/approxdoccount", jwtAuthMiddleware(legacy.ApproxDocCount))
		// index management (stop-the-world)
		router.GET("/index", jwtAuthMiddleware(legacy.Index))
		router.GET("/indexes", jwtAuthMiddleware(legacy.Indexes))
		router.GET("/unindex", jwtAuthMiddleware(legacy.Unindex))
		// misc (stop-the-world)
		router.GET("/shutdown", jwtAuthMiddleware(legacy.Shutdown))
		router.GET("/dump", jwtAuthMiddleware(legacy.Dump))
	}

	// New API
	// These endpoints are always available and do not require authentication
	router.GET("/", middlewares.StandardResponse(Welcome))
	router.GET("/version", middlewares.StandardResponse(legacy.Version))
	router.GET("/memstats", middlewares.StandardResponse(legacy.MemStats))

	// collection management (stop-the-world)
	router.POST("/collection/:collection_name", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Create)))
	router.PUT("/collection/:collection_name/rename/:new_collection_name", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Rename)))
	router.DELETE("/collection/:collection_name", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Drop)))
	router.GET("/collections", jwtAuthMiddleware(middlewares.StandardResponse(legacy.All)))
	router.POST("/collection/:collection_name/scrub", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Scrub)))
	router.POST("/sync", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Sync)))
	// query
	router.POST("/collection/:collection_name/query", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Query)))
	router.POST("/collection/:collection_name/count", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Count)))
	// document management
	router.POST("/collection/:collection_name/doc",jwtAuthMiddleware(middlewares.StandardResponse(legacy.Insert)))
	router.GET("/collection/:collection_name/doc/:id", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Get)))
	router.PUT("/collection/:collection_name/doc/:id", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Update)))
	router.DELETE("/collection/:collection_name/doc/:id", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Delete)))
	router.GET("/collection/:collection_name/page/:page/of/:total", jwtAuthMiddleware(middlewares.StandardResponse(legacy.GetPage)))
	// TODO: Review if it will make more sense for it to be just /collection/:collection_name/count
	router.GET("/collection/:collection_name/count/approx", jwtAuthMiddleware(middlewares.StandardResponse(legacy.ApproxDocCount)))
	// index management (stop-the-world)
	router.POST("/collection/:collection_name/index", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Index)))
	router.DELETE("/collection/:collection_name/index", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Unindex)))
	router.GET("/collection/:collection_name/indexes", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Indexes)))
	// misc (stop-the-world)
	router.POST("/shutdown", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Shutdown)))
	router.POST("/dump", jwtAuthMiddleware(middlewares.StandardResponse(legacy.Dump)))

	iface := "all interfaces"
	if bind != "" {
		iface = bind
	}

	if tlsCrt != "" {
		tdlog.Noticef("Will listen on %s (HTTPS), port %d.", iface, port)
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", bind, port), tlsCrt, tlsKey, nil); err != nil {
			tdlog.Panicf("Failed to start HTTPS service - %s", err)
		}
	} else {
		tdlog.Noticef("Will listen on %s (HTTP), port %d.", iface, port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), nil)
	}
}

// Greet user with a welcome message.
func Welcome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Path != "/" {
		http.Error(w, "Invalid API endpoint", 404)
		return
	}
	w.Write([]byte("Welcome to tiedot"))
}
