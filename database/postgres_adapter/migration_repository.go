package postgresadapter

import (
	"github.com/mhogar/migrationrunner"
)

// PostgresMigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the postgres db.
type PostgresMigrationRepository struct {
	DB *PostgresDB
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo PostgresMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{}
}
