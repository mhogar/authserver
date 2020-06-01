package helpers

import "github.com/google/uuid"

// TokenHelper is an interface for encapsulating helper methods for creating and parsing tokens.
type TokenHelper interface {
	// CreateSessionToken creates a session token to be used for authenticating requests.
	CreateSessionToken(sID uuid.UUID) (string, error)
}
