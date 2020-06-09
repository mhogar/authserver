package models

import (
	"github.com/google/uuid"
)

// Scope ValidateError statuses.
const (
	ValidateScopeValid       = iota
	ValidateScopeInvalidID   = iota
	ValidateScopeInvalidName = iota
)

// Scope represents the scope model.
type Scope struct {
	ID   uuid.UUID
	Name string
}

// ScopeCRUD is an interface for performing CRUD operations on a scope.
type ScopeCRUD interface {
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
// Returns a ValidateError indicating its result.
func (s *Scope) Validate() ValidateError {
	if s.ID == uuid.Nil {
		return CreateValidateError(ValidateScopeInvalidID, "id cannot be nil")
	}

	if s.Name == "" {
		return CreateValidateError(ValidateScopeInvalidName, "name cannot be empty")
	}

	return ValidateError{ValidateScopeValid, nil}
}
