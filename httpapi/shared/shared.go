package shared

import (
	"github.com/HouzuoGuo/tiedot/db"
)

var databaseInstance *db.DB // HTTP API endpoints operate on this database

// InitializeDatabaseInstance Initializes Database Instance
func InitializeDatabaseInstance(databasePath string) error {
	if databaseInstance == nil {
		db, err := db.OpenDB(databasePath)
		if err != nil {
			return err
		}
		databaseInstance = db
	}
	return nil
}

// GetDatabaseInstance Gets the Shared Database Instance
func GetDatabaseInstance() *db.DB {
	return databaseInstance
}
