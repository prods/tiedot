package api

import "github.com/HouzuoGuo/tiedot/db"

var databseInstance *db.DB // HTTP API endpoints operate on this database

// Start Starts API Server
func Start(dir string, port int, tlsCrt, tlsKey, jwtPubKey, jwtPrivateKey, bind, authToken string) {
	var err error
	databseInstance, err = db.OpenDB(dir)
	if err != nil {
		panic(err)
	}
}
