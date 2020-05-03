package service

import (
	"github.com/devingen/api-core/database"
	"github.com/devingen/data-api/service/database-service"
)

// NewDatabaseService generates new DatabaseService
func NewDatabaseService(database *database.Database) *database_service.DatabaseService {
	return &database_service.DatabaseService{
		Database: database,
	}
}
