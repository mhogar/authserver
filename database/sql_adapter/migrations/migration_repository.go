package migrations

import (
	sqladapter "authserver/database/sql_adapter"

	"github.com/mhogar/migrationrunner"
)

// SQLMigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the sql adapter.
type SQLMigrationRepository struct {
	Adapter *sqladapter.SQLAdapter
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo SQLMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		m20200611224101{Adapter: repo.Adapter},
	}
}
