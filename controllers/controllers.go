package controllers

import (
	requesterror "authserver/common/request_error"
	"authserver/models"

	"github.com/google/uuid"
)

// Controllers encapsulates all other controller interfaces.
type Controllers interface {
	UserController
	TokenController
}

// UserControllerCRUD encapsulates the CRUD operations required by the UserController.
type UserControllerCRUD interface {
	models.UserCRUD
	models.AccessTokenCRUD
}

// UserController provides workflows for user related operations.
type UserController interface {
	// CreateUser creates a new user with the given username and password.
	CreateUser(CRUD UserControllerCRUD, username string, password string) (*models.User, requesterror.RequestError)

	// DeleteUser deletes the given user.
	DeleteUser(CRUD UserControllerCRUD, user *models.User) requesterror.RequestError

	// UpdateUserPassword updates the given user's password.
	UpdateUserPassword(CRUD UserControllerCRUD, user *models.User, oldPassword string, newPassword string) requesterror.RequestError
}

// TokenControllerCRUD encapsulates the CRUD operations required by the TokenController.
type TokenControllerCRUD interface {
	models.UserCRUD
	models.ClientCRUD
	models.ScopeCRUD
	models.AccessTokenCRUD
}

// TokenController provides workflows for access token related operations.
type TokenController interface {
	// CreateTokenFromPassword creates a new access token, authenticating using a password.
	CreateTokenFromPassword(CRUD TokenControllerCRUD, username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, requesterror.OAuthRequestError)

	// DeleteToken deletes the access token.
	DeleteToken(CRUD TokenControllerCRUD, token *models.AccessToken) requesterror.RequestError
}

// Controls encapsulates all other control structs.
type Controls struct {
	UserControl
	TokenControl
}
