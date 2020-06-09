package models

import "github.com/google/uuid"

// AccessToken ValidateError statuses.
const (
	ValidateAccessTokenValid           = iota
	ValidateAccessTokenInvalidID       = iota
	ValidateAccessTokenInvalidUserID   = iota
	ValidateAccessTokenInvalidClientID = iota
	ValidateAccessTokenInvalidScopeID  = iota
)

// AccessToken represents the access token model.
type AccessToken struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ClientID uuid.UUID
	ScopeID  uuid.UUID
}

// AccessTokenCRUD is an interface for performing CRUD operations on an access token.
type AccessTokenCRUD interface {
	// SaveAccessToken saves the access token and returns any errors.
	SaveAccessToken(token *AccessToken) error

	// GetAccessTokenByID fetches the access token associated with the id.
	// If no tokens are found, returns nil token. Also returns any errors.
	GetAccessTokenByID(ID uuid.UUID) (*AccessToken, error)

	// DeleteAccessToken deletes the token and returns any errors.
	DeleteAccessToken(token *AccessToken) error
}

// CreateNewAccessToken creates a access token model with a new id and the provided fields.
func CreateNewAccessToken(userID uuid.UUID, clientID uuid.UUID, scopeID uuid.UUID) *AccessToken {
	return &AccessToken{
		ID:       uuid.New(),
		UserID:   userID,
		ClientID: clientID,
		ScopeID:  scopeID,
	}
}

// Validate validates the access token model has valid fields.
// Returns a ValidateError indicating its result.
func (tk *AccessToken) Validate() ValidateError {
	if tk.ID == uuid.Nil {
		return CreateValidateError(ValidateAccessTokenInvalidID, "id cannot be nil")
	}

	if tk.UserID == uuid.Nil {
		return CreateValidateError(ValidateAccessTokenInvalidUserID, "user id cannot be nil")
	}

	if tk.ClientID == uuid.Nil {
		return CreateValidateError(ValidateAccessTokenInvalidClientID, "client id cannot be nil")
	}

	if tk.ScopeID == uuid.Nil {
		return CreateValidateError(ValidateAccessTokenInvalidScopeID, "scope id cannot be nil")
	}

	return ValidateError{ValidateAccessTokenValid, nil}
}
