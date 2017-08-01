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
	"io/ioutil"
	"net/http"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"strings"
)

var (
	HttpDB *db.DB // HTTP API endpoints operate on this database
)

// IsNewAPIRoute Determines if the provided Route is part of the new API.
func IsNewAPIRoute(r *http.Request) bool {
	// This is a temporary solution in order to reuse the original route handlers.
	// It should be gone once the old API is deprecated.
	return strings.HasPrefix(r.URL.Path, "/collection") ||
		(r.Method == "POST" && (r.URL.Path == "/sync" ||
			r.URL.Path == "/shutdown" ||
			r.URL.Path == "/dump"))
}

// Store form parameter value of specified key to *val and return true; if key does not exist, set HTTP status 400 and return false.
func Require(w http.ResponseWriter, r *http.Request, key string, val *string) bool {
	*val = r.FormValue(key)
	if *val == "" {
		http.Error(w, fmt.Sprintf("Please pass POST/PUT/GET parameter value of '%s'.", key), 400)
		return false
	}
	return true
}

// Start HTTP server and block until the server shuts down. Panic on error.
func Start(dir string, port int, tlsCrt, tlsKey, jwtPubKey, jwtPrivateKey, bind, authToken string, backwardsCompatibleAPI bool) {
	var err error
	HttpDB, err = db.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	// These endpoints are always available and do not require authentication
	router.GET("/", StandardAPIResponder(Welcome))
	router.GET("/version", StandardAPIResponder(Version))
	router.GET("/memstats", StandardAPIResponder(MemStats))

	// Install API endpoint handlers that may require authorization
	var authWrap func(httprouter.Handle) httprouter.Handle
	if authToken != "" {
		tdlog.Noticef("API endpoints now require the pre-shared token in Authorization header.")
		authWrap = func(originalHandler httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				if "token "+authToken != r.Header.Get("Authorization") {
					http.Error(w, "", http.StatusUnauthorized)
					return
				}
				originalHandler(w, r, ps)
			}
		}
	} else if jwtPubKey != "" && jwtPrivateKey != "" {
		tdlog.Noticef("API endpoints now require JWT in Authorization header.")
		var publicKeyContent, privateKeyContent []byte
		if publicKeyContent, err = ioutil.ReadFile(jwtPubKey); err != nil {
			panic(err)
		} else if publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyContent); err != nil {
			panic(err)
		} else if privateKeyContent, err = ioutil.ReadFile(jwtPrivateKey); err != nil {
			panic(err)
		} else if privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyContent); err != nil {
			panic(err)
		}
		jwtInitSetup()
		authWrap = jwtWrap
		// does not require JWT auth
		http.HandleFunc("/getjwt", getJWT)
		http.HandleFunc("/checkjwt", checkJWT)
	} else {
		tdlog.Noticef("API endpoints do not require Authorization header.")
		authWrap = func(originalHandler httprouter.Handle) httprouter.Handle {
			return originalHandler
		}
	}

	// Legacy API Routes
	if backwardsCompatibleAPI {
		// collection management (stop-the-world)
		router.GET("/create", authWrap(Create))
		router.GET("/rename", authWrap(Rename))
		router.GET("/drop", authWrap(Drop))
		router.GET("/all", authWrap(All))
		router.GET("/scrub", authWrap(Scrub))
		router.GET("/sync", authWrap(Sync))
		// query
		router.GET("/query", authWrap(Query))
		router.GET("/count", authWrap(Count))
		// document management
		router.GET("/insert", authWrap(Insert))
		router.GET("/get", authWrap(Get))
		router.GET("/getpage", authWrap(GetPage))
		router.GET("/update", authWrap(Update))
		router.GET("/delete", authWrap(Delete))
		router.GET("/approxdoccount", authWrap(ApproxDocCount))
		// index management (stop-the-world)
		router.GET("/index", authWrap(Index))
		router.GET("/indexes", authWrap(Indexes))
		router.GET("/unindex", authWrap(Unindex))

		// misc (stop-the-world)
		router.GET("/shutdown", authWrap(Shutdown))
		router.GET("/dump", authWrap(Dump))
	}

	// collection management (stop-the-world)
	router.POST("/collection/:collection_name", authWrap(StandardAPIResponder(Create)))
	router.PUT("/collection/:collection_name/rename/:new_collection_name", authWrap(StandardAPIResponder(Rename)))
	router.DELETE("/collection/:collection_name", authWrap(StandardAPIResponder(Drop)))
	router.GET("/collections", authWrap(StandardAPIResponder(All)))
	router.POST("/collection/:collection_name/scrub", authWrap(StandardAPIResponder(Scrub)))
	router.POST("/sync", authWrap(StandardAPIResponder(Sync)))
	// query
	router.POST("/collection/:collection_name/query", authWrap(StandardAPIResponder(Query)))
	router.POST("/collection/:collection_name/count", authWrap(StandardAPIResponder(Count)))
	// document management
	router.POST("/collection/:collection_name/doc", authWrap(StandardAPIResponder(Insert)))
	router.GET("/collection/:collection_name/doc/:id", authWrap(StandardAPIResponder(Get)))
	router.PUT("/collection/:collection_name/doc/:id", authWrap(StandardAPIResponder(Update)))
	router.DELETE("/collection/:collection_name/doc/:id", authWrap(StandardAPIResponder(Delete)))
	router.GET("/collection/:collection_name/page/:page/of/:total", authWrap(StandardAPIResponder(GetPage)))
	// TODO: Review if it will make more sense for it to be just /collection/:collection_name/count
	router.GET("/collection/:collection_name/count/approx", authWrap(StandardAPIResponder(ApproxDocCount)))
	// index management (stop-the-world)
	router.POST("/collection/:collection_name/index", authWrap(StandardAPIResponder(Index)))
	router.DELETE("/collection/:collection_name/index", authWrap(StandardAPIResponder(Unindex)))
	router.GET("/collection/:collection_name/indexes", authWrap(StandardAPIResponder(Indexes)))
	// misc (stop-the-world)
	router.POST("/shutdown", authWrap(StandardAPIResponder(Shutdown)))
	router.POST("/dump", authWrap(StandardAPIResponder(Dump)))

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
		http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), router)
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
