package controllers

import (
	requesterror "authserver/common/request_error"
	"authserver/models"

	"github.com/google/uuid"
)

type Controllers interface {
	UserController
	TokenController
}

type UserController interface {
	// CreateUser creates a new user with the given username and password
	CreateUser(username string, password string) (*models.User, requesterror.RequestError)

	// DeleteUser deletes the user with the given id
	DeleteUser(id uuid.UUID) requesterror.RequestError

	// UpdateUserPassword updates the given user's password
	UpdateUserPassword(user *models.User, oldPassword string, newPassword string) requesterror.RequestError
}

type TokenController interface {
	// CreateTokenFromPassword creates a new access token, authenticating using a password
	CreateTokenFromPassword(username string, password string, clientID string, scopeName string) (*models.AccessToken, requesterror.OAuthRequestError)

	// DeleteToken deletes the access token
	DeleteToken(token *models.AccessToken) requesterror.RequestError
}

type BaseControllers struct {
	UserControl
	TokenControl
}
