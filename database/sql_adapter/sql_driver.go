package sqladapter

// SQLDriver is an interface for encapsulating methods specific to each sql driver.
type SQLDriver interface {
	SQLScriptRepository

	// GetDriverName returns the name for the driver.
	GetDriverName() string
}

// SQLScriptRepository is an interface for fetching sql scripts.
type SQLScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string

	CreateUserTableScript() string
	DropUserTableScript() string
	SaveUserScript() string
	GetUserByIdScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	DeleteUserScript() string

	CreateClientTableScript() string
	DropClientTableScript() string
	SaveClientScript() string
	GetClientByIdScript() string
}
