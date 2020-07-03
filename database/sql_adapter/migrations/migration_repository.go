package migrations

import (
	sqladapter "authserver/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

// MigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the sql db.
type MigrationRepository struct {
	DB *sqladapter.SQLDB
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo MigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		m20200628151601{DB: repo.DB},
	}
}
