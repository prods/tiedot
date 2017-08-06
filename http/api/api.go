package api

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	"github.com/HouzuoGuo/tiedot/tdlog"
	"net/http"
	"fmt"
	"github.com/HouzuoGuo/tiedot/http/api/modules"
)

var databaseInstance *db.DB // HTTP API endpoints operate on this database

// Start Starts API Server
func Start(dir string, port int, tlsCrt, tlsKey, jwtPubKey, jwtPrivateKey, bind, authToken string) {

	// Database Instance Initialization
	var err error
	databaseInstance, err = db.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	// Router Initialization
	router := httprouter.New()

	// Register API Modules

	// Miscellaneous API Module
	miscellaneousAPIModule := modules.NewMiscellaneousAPIModule(databaseInstance)
	modules.Mount(router, miscellaneousAPIModule)

	// Indexes API Module
	indexesAPIModule := modules.NewIndexesAPIModule(databaseInstance)
	modules.Mount(router, indexesAPIModule)


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
