package database_service

import (
	"github.com/devingen/api-core/database"
)

// DatabaseService implements DataService interface with database connection
type DatabaseService struct {
	Database *database.Database
}
