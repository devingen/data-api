package mongo_data_service

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/data-api/data-service"
)

// MongoDataService implements DataService interface with database connection
type MongoDataService struct {
	Database *database.Database
}

// New generates new DatabaseService
func New(database *database.Database) data_service.IDataService {
	return MongoDataService{
		Database: database,
	}
}
