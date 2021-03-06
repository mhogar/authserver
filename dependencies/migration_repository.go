package dependencies

import (
	sqladapter "authserver/database/sql_adapter"
	"authserver/database/sql_adapter/migrations"
	"sync"

	"github.com/mhogar/migrationrunner"
)

var createMigrationRepositoryOnce sync.Once
var migrationRepository migrationrunner.MigrationRepository

// ResolveMigrationRepository resolves the MigrationRepository dependency.
// Only the first call to this function will create a new MigrationRepository, after which it will be retrieved from memory.
func ResolveMigrationRepository() migrationrunner.MigrationRepository {
	createMigrationRepositoryOnce.Do(func() {
		migrationRepository = migrations.MigrationRepository{
			DB: ResolveDatabase().(*sqladapter.SQLDB),
		}
	})
	return migrationRepository
}
