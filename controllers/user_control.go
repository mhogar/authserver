package controllers

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	requesterror "authserver/common/request_error"
	"authserver/helpers"
	commonhelpers "authserver/helpers/common"
	"authserver/models"
)

// UserControl handles requests to "/user" endpoints
type UserControl struct {
	CRUD interface {
		models.UserCRUD
		models.AccessTokenCRUD
	}
	PasswordHasher            helpers.PasswordHasher
	PasswordCriteriaValidator helpers.PasswordCriteriaValidator
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
		log.Println(commonhelpers.ChainError("error getting user by username", err))
		return nil, requesterror.InternalError()
	}
	if otherUser != nil {
		return nil, requesterror.ClientError("username is already in use")
	}

	//validate password meets criteria
	vperr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(password)
	if vperr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(commonhelpers.ChainError("error validating password criteria", vperr))
		return nil, requesterror.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	user.PasswordHash, err = c.PasswordHasher.HashPassword(password)
	if err != nil {
		log.Println(commonhelpers.ChainError("error generating password hash", err))
		return nil, requesterror.InternalError()
	}

	//save the user
	err = c.CRUD.SaveUser(user)
	if err != nil {
		log.Println(commonhelpers.ChainError("error saving user", err))
		return nil, requesterror.InternalError()
	}

	return user, requesterror.NoError()
}

// DeleteUser deletes the user with the given id
func (c UserControl) DeleteUser(id uuid.UUID) requesterror.RequestError {
	//get the user
	user, err := c.CRUD.GetUserByID(id)
	if err != nil {
		log.Println(commonhelpers.ChainError("error fetching user by id", err))
		return requesterror.InternalError()
	}

	if user == nil {
		return requesterror.ClientError("user not found")
	}

	//delete the user
	err = c.CRUD.DeleteUser(user)
	if err != nil {
		log.Println(commonhelpers.ChainError("error deleting user", err))
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
		log.Println(commonhelpers.ChainError("error comparing password hashes", err))
		return requesterror.ClientError("old password is invalid")
	}

	//validate new password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(newPassword)
	if verr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(commonhelpers.ChainError("error validating password criteria", verr))
		return requesterror.ClientError("password does not meet minimum criteria")
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(newPassword)
	if err != nil {
		log.Println(commonhelpers.ChainError("error generating password hash", err))
		return requesterror.InternalError()
	}

	//update the user
	user.PasswordHash = hash
	err = c.CRUD.UpdateUser(user)
	if err != nil {
		log.Println(commonhelpers.ChainError("error updating user", err))
		return requesterror.InternalError()
	}

	//return success
	return requesterror.NoError()
}
