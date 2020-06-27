package postgresadapter

import (
	"authserver/database/postgres_adapter/scripts"
)

type PostgresScriptRepository struct{}

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
