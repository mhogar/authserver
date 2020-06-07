package models

import (
	"github.com/google/uuid"
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
