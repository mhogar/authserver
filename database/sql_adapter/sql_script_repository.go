package sqladapter

// SQLScriptRepository is an interface for fetching sql scripts.
type SQLScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string
	CreateUserTableScript() string
	SaveUserScript() string
	GetUserByIDScript() string
}
