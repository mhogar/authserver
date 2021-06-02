package controllers

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	passwordhelpers "authserver/controllers/password_helpers"
	"authserver/models"
	"log"

	"github.com/google/uuid"
)

// TokenControl handles requests to "/token" endpoints
type TokenControl struct {
	CRUD interface {
		models.UserCRUD
		models.ClientCRUD
		models.ScopeCRUD
		models.AccessTokenCRUD
	}
	passwordhelpers.PasswordHasher
}

// PostToken handles POST requests to "/token"
func (c TokenControl) CreateTokenFromPassword(username string, password string, clientID uuid.UUID, scopeName string) (*models.AccessToken, requesterror.OAuthRequestError) {
	//get the client
	client, rerr := parseClient(c.CRUD, clientID)
	if rerr.Type != requesterror.ErrorTypeNone {
		return nil, rerr
	}

	//get the scope
	scope, rerr := parseScope(c.CRUD, scopeName)
	if rerr.Type != requesterror.ErrorTypeNone {
		return nil, rerr
	}

	//get the user
	user, err := c.CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, requesterror.OAuthInternalError()
	}

	//check if user was found
	if user == nil {
		return nil, requesterror.OAuthClientError("invalid_grant", "invalid username and/or password")
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, password)
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return nil, requesterror.OAuthClientError("invalid_grant", "invalid username and/or password")
	}

	//create a new access token
	token := models.CreateNewAccessToken(user, client, scope)

	//save the token
	err = c.CRUD.SaveAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error saving access token", err))
		return nil, requesterror.OAuthInternalError()
	}

	return token, requesterror.OAuthNoError()
}

// DeleteToken deletes the access token.
func (c TokenControl) DeleteToken(token *models.AccessToken) requesterror.RequestError {
	//delete the token
	err := c.CRUD.DeleteAccessToken(token)
	if err != nil {
		log.Println(common.ChainError("error deleting access token", err))
		return requesterror.InternalError()
	}

	//return success
	return requesterror.NoError()
}
