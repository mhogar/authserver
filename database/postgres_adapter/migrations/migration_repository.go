package migrations

import (
	postgresadpater "authserver/database/postgres_adapter"

	"github.com/mhogar/migrationrunner"
)

// PostgresMigrationRepository is an implementation of the MigrationRepository interface that fetches migrations for the postgres db.
type PostgresMigrationRepository struct {
	DB *postgresadpater.PostgresDB
}

// GetMigrations returns a slice of Migrations that need to be run on the sql database.
func (repo PostgresMigrationRepository) GetMigrations() []migrationrunner.Migration {
	return []migrationrunner.Migration{
		m20200628151601{DB: repo.DB},
	}
}
