package sqladapter

// SQLScriptRepository is an interface for encapsulating other sql script repository.
type SQLScriptRepository interface {
	ClientScriptRepository
	MigrationScriptRepository
	UserScriptRepository
}

// ClientScriptRepository is an interface for fetching client sql scripts.
type ClientScriptRepository interface {
	CreateClientTableScript() string
	DropClientTableScript() string
	SaveClientScript() string
	GetClientByIdScript() string
}

// MigrationScriptRepository is an interface for fetching migration sql scripts.
type MigrationScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string
}

// UserScriptRepository is an interface for fetching user sql scripts.
type UserScriptRepository interface {
	CreateUserTableScript() string
	DropUserTableScript() string
	SaveUserScript() string
	GetUserByIdScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	DeleteUserScript() string
}
