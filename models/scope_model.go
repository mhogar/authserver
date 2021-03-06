package models

import (
	"github.com/google/uuid"
)

// Scope ValidateError statuses.
const (
	ValidateScopeValid       = 0x0
	ValidateScopeNilID       = 0x1
	ValidateScopeEmptyName   = 0x2
	ValidateScopeNameTooLong = 0x4
)

// ScopeNameMaxLength is the max length a scope's name can be.
const ScopeNameMaxLength = 15

// Scope represents the scope model.
type Scope struct {
	ID   uuid.UUID
	Name string
}

// ScopeCRUD is an interface for performing CRUD operations on a scope.
type ScopeCRUD interface {
	// SaveScope saves the scope and returns any errors.
	SaveScope(scope *Scope) error

	// GetScopeByName fetches the scope with the matching name.
	// If no scopes are found, returns nil scope. Also returns any errors.
	GetScopeByName(name string) (*Scope, error)
}

// CreateNewScope creates a Scope model with a new id and the provided fields.
func CreateNewScope(name string) *Scope {
	return &Scope{
		ID:   uuid.New(),
		Name: name,
	}
}

// Validate validates the client model has valid fields.
// Returns an int indicating which fields are invalid.
func (s *Scope) Validate() int {
	code := ValidateScopeValid

	if s.ID == uuid.Nil {
		code |= ValidateScopeNilID
	}

	if s.Name == "" {
		code |= ValidateScopeEmptyName
	} else if len(s.Name) > ScopeNameMaxLength {
		code |= ValidateScopeNameTooLong
	}

	return code
}
