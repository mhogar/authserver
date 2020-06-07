package models

import (
	"github.com/google/uuid"
)

// Scope represents the scope model.
type Scope struct {
	ID   uuid.UUID
	Name string
}

// CreateNewScope creates a Scope model with a new id and the provided fields.
func CreateNewScope(name string) *Scope {
	return &Scope{
		ID:   uuid.New(),
		Name: name,
	}
}
