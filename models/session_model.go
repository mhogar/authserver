package models

import (
	"github.com/google/uuid"
)

// Session ValidateError statuses.
const (
	ValidateSessionValid         = iota
	ValidateSessionInvalidID     = iota
	ValidateSessionInvalidUserID = iota
)

// Session represents the session model.
type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// CreateNewSession creates a session model with a new id and the provided fields.
func CreateNewSession(userID uuid.UUID) *Session {
	return &Session{
		ID:     uuid.New(),
		UserID: userID,
	}
}

// CreateValidateSessionValid creates a ValidateError with status ValidateSessionValid and nil error.
func CreateValidateSessionValid() ValidateError {
	return ValidateError{ValidateSessionValid, nil}
}

// Validate validates that the session model has valid fields.
// Returns a ValidateError indicating its result.
func (s *Session) Validate() ValidateError {
	if s.ID == uuid.Nil {
		return CreateValidateError(ValidateSessionInvalidID, "id cannot be nil")
	}

	if s.UserID == uuid.Nil {
		return CreateValidateError(ValidateSessionInvalidUserID, "user id cannot be nil")
	}

	return CreateValidateSessionValid()
}
