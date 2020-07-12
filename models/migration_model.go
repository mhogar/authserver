package models

import (
	"regexp"

	"github.com/mhogar/migrationrunner"
)

// Migration ValidateError statuses.
const (
	ValidateMigrationValid            = iota
	ValidateMigrationInvalidTimestamp = iota
)

// Migration represents the migration model.
type Migration struct {
	Timestamp string
}

// MigrationCRUD is an interface for performing CRUD operations on a migration.
type MigrationCRUD interface {
	migrationrunner.MigrationCRUD

	// GetMigrationByTimestamp fetches the migration with the matching timestamp.
	// If no migrations are found, returns nil migration. Also returns any errors.
	GetMigrationByTimestamp(timestamp string) (*Migration, error)
}

// CreateNewMigration creates a new migration with the given timestamp.
func CreateNewMigration(timestamp string) *Migration {
	return &Migration{
		Timestamp: timestamp,
	}
}

// Validate validates the migration is a valid migration model.
// Returns a ValidateError indicating its result.
func (m Migration) Validate() ValidateError {
	matched, _ := regexp.MatchString(`^\d{14}$`, m.Timestamp)
	if !matched {
		return CreateValidateError(ValidateMigrationInvalidTimestamp, "timestamp is in invalid format")
	}

	return ValidateError{ValidateMigrationValid, nil}
}
