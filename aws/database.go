package aws

import (
	"github.com/devingen/api-core/database"
	"log"
)

var db *database.Database

func GetDatabase() *database.Database {
	if db == nil {
		var err error
		db, err = database.NewDatabase()
		if err != nil {
			log.Fatalf("Database connection failed %s", err.Error())
		}
	} else if !db.IsConnected() {
		err := db.ConnectWithEnvironment()
		if err != nil {
			log.Fatalf("Database connection failed %s", err.Error())
		}
	} else {
		log.Println("Database connection exists")
	}
	return db
}
