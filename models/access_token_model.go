package models

import "github.com/google/uuid"

// AccessToken ValidateError statuses.
const (
	ValidateAccessTokenValid         = iota
	ValidateAccessTokenNilID         = iota
	ValidateAccessTokenNilUser       = iota
	ValidateAccessTokenInvalidUser   = iota
	ValidateAccessTokenNilClient     = iota
	ValidateAccessTokenInvalidClient = iota
	ValidateAccessTokenNilScope      = iota
	ValidateAccessTokenInvalidScope  = iota
)

// AccessToken represents the access token model.
type AccessToken struct {
	ID     uuid.UUID
	User   *User
	Client *Client
	Scope  *Scope
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
func CreateNewAccessToken(user *User, client *Client, scope *Scope) *AccessToken {
	return &AccessToken{
		ID:     uuid.New(),
		User:   user,
		Client: client,
		Scope:  scope,
	}
}

// Validate validates the access token model has valid fields.
// Returns a ValidateError indicating its result.
func (tk *AccessToken) Validate() ValidateError {
	if tk.ID == uuid.Nil {
		return CreateValidateError(ValidateAccessTokenNilID, "id cannot be nil")
	}

	if tk.User == nil {
		return CreateValidateError(ValidateAccessTokenNilUser, "user cannot be nil")
	}

	verr := tk.User.Validate()
	if verr.Status != ValidateUserValid {
		return CreateValidateError(ValidateAccessTokenInvalidUser, "invalid user: "+verr.Error())
	}

	if tk.Client == nil {
		return CreateValidateError(ValidateAccessTokenNilClient, "client cannot be nil")
	}

	verr = tk.Client.Validate()
	if verr.Status != ValidateClientValid {
		return CreateValidateError(ValidateAccessTokenInvalidClient, "invalid client: "+verr.Error())
	}

	if tk.Scope == nil {
		return CreateValidateError(ValidateAccessTokenNilScope, "scope cannot be nil")
	}

	verr = tk.Scope.Validate()
	if verr.Status != ValidateScopeValid {
		return CreateValidateError(ValidateAccessTokenInvalidScope, "invalid scope: "+verr.Error())
	}

	return ValidateError{ValidateAccessTokenValid, nil}
}
