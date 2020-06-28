package dependencies

import (
	postgresadapter "authserver/database/postgres_adapter"
	"sync"

	"github.com/mhogar/migrationrunner"
)

var createMigrationRepositoryOnce sync.Once
var migrationRepository migrationrunner.MigrationRepository

// ResolveMigrationRepository resolves the MigrationRepository dependency.
// Only the first call to this function will create a new MigrationRepository, after which it will be retrieved from the cache.
func ResolveMigrationRepository() migrationrunner.MigrationRepository {
	createMigrationRepositoryOnce.Do(func() {
		migrationRepository = postgresadapter.PostgresMigrationRepository{
			DB: ResolveDatabase().(*postgresadapter.PostgresDB),
		}
	})
	return migrationRepository
}
