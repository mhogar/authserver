package controllers

import (
	"fmt"
	"log"

	"authserver/common"
	requesterror "authserver/common/request_error"
	passwordhelpers "authserver/controllers/password_helpers"
	"authserver/models"
)

// UserControl handles requests to "/user" endpoints
type UserControl struct {
	CRUD interface {
		models.UserCRUD
		models.AccessTokenCRUD
	}
	PasswordHasher            passwordhelpers.PasswordHasher
	PasswordCriteriaValidator passwordhelpers.PasswordCriteriaValidator
}

// CreateUser creates a new user with the given username and password
func (c UserControl) CreateUser(username string, password string) (*models.User, requesterror.RequestError) {
	//create the user model
	user := models.CreateNewUser(username, nil)

	//validate the username
	verr := user.Validate()
	if verr&models.ValidateUserEmptyUsername != 0 {
		return nil, requesterror.ClientError("username cannot be empty")
	} else if verr&models.ValidateUserUsernameTooLong != 0 {
		return nil, requesterror.ClientError(fmt.Sprint("username cannot be longer than ", models.UserUsernameMaxLength, " characters"))
	}

	//validate username is unique
	otherUser, err := c.CRUD.GetUserByUsername(username)
	if err != nil {
		log.Println(common.ChainError("error getting user by username", err))
		return nil, requesterror.InternalError()
	}
	if otherUser != nil {
		return nil, requesterror.ClientError("error creating user")
	}

	//validate password meets criteria
	vperr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(password)
	if vperr.Status != passwordhelpers.ValidatePasswordCriteriaValid {
		log.Println(common.ChainError("error validating password criteria", vperr))
		return nil, requesterror.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	user.PasswordHash, err = c.PasswordHasher.HashPassword(password)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return nil, requesterror.InternalError()
	}

	//save the user
	err = c.CRUD.SaveUser(user)
	if err != nil {
		log.Println(common.ChainError("error saving user", err))
		return nil, requesterror.InternalError()
	}

	return user, requesterror.NoError()
}

// DeleteUser deletes the user with the given id
func (c UserControl) DeleteUser(user *models.User) requesterror.RequestError {
	//delete the user
	err := c.CRUD.DeleteUser(user)
	if err != nil {
		log.Println(common.ChainError("error deleting user", err))
		return requesterror.InternalError()
	}

	//return success
	return requesterror.NoError()
}

// UpdateUserPassword updates the given user's password
func (c UserControl) UpdateUserPassword(user *models.User, oldPassword string, newPassword string) requesterror.RequestError {
	//validate old password
	err := c.PasswordHasher.ComparePasswords(user.PasswordHash, oldPassword)
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return requesterror.ClientError("old password is invalid")
	}

	//validate new password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(newPassword)
	if verr.Status != passwordhelpers.ValidatePasswordCriteriaValid {
		log.Println(common.ChainError("error validating password criteria", verr))
		return requesterror.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(newPassword)
	if err != nil {
		log.Println(common.ChainError("error generating password hash", err))
		return requesterror.InternalError()
	}

	//update the user
	user.PasswordHash = hash
	err = c.CRUD.UpdateUser(user)
	if err != nil {
		log.Println(common.ChainError("error updating user", err))
		return requesterror.InternalError()
	}

	//TODO: delete all other user access tokens (requires transaction)

	//return success
	return requesterror.NoError()
}
