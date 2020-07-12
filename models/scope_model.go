package models

import (
	"fmt"

	"github.com/google/uuid"
)

// Scope ValidateError statuses.
const (
	ValidateScopeValid       = iota
	ValidateScopeNilID       = iota
	ValidateScopeEmptyName   = iota
	ValidateScopeNameTooLong = iota
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
// Returns a ValidateError indicating its result.
func (s *Scope) Validate() ValidateError {
	if s.ID == uuid.Nil {
		return CreateValidateError(ValidateScopeNilID, "id cannot be nil")
	}

	if s.Name == "" {
		return CreateValidateError(ValidateScopeEmptyName, "name cannot be empty")
	}

	if len(s.Name) > ScopeNameMaxLength {
		return CreateValidateError(ValidateScopeNameTooLong, fmt.Sprint("name cannot be longer than", ScopeNameMaxLength, "characters"))
	}

	return ValidateError{ValidateScopeValid, nil}
}
