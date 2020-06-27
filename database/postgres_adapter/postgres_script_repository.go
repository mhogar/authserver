package postgresadapter

import (
	"authserver/database/postgres_adapter/scripts"
)

// PostgresScriptRepository is an implementation of the sql script respository interface for postgres scripts.
type PostgresScriptRepository struct{}

// GetSQLScript gets the sql script for the given key.
func (PostgresScriptRepository) GetSQLScript(key string) string {
	switch key {
	case "CreateMigrationTable":
		return scripts.GetCreateMigrationTableScript()
	case "SaveMigration":
		return scripts.GetSaveMigrationScript()
	case "GetMigrationByTimestamp":
		return scripts.GetGetMigrationByTimestampScript()
	case "GetLatestTimestamp":
		return scripts.GetGetLatestTimestampScript()
	case "DeleteMigrationByTimestamp":
		return scripts.GetDeleteMigrationByTimestampScript()
	default:
		return ""
	}
}
