package controllers

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"authserver/database"
	"authserver/helpers"
	"authserver/models"

	"github.com/julienschmidt/httprouter"
)

// UserController handles requests to "/user" endpoints
type UserController struct {
	UserCRUD                  database.UserCRUD
	PasswordHasher            helpers.PasswordHasher
	PasswordCriteriaValidator helpers.PasswordCriteriaValidator
}

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (c UserController) PostUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostUserBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PostUser request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//validate the body fields
	if body.Username == "" || body.Password == "" {
		sendErrorResponse(w, http.StatusBadRequest, "username and password cannot be empty")
		return
	}

	//validate username is unique
	otherUser, err := c.UserCRUD.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by username", err))
		sendInternalErrorResponse(w)
		return
	}

	if otherUser != nil {
		sendErrorResponse(w, http.StatusBadRequest, "an user with that username already exists")
		return
	}

	//validate password meets criteria
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(body.Password)
	if verr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(helpers.ChainError("error validating password criteria", verr))
		sendErrorResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(body.Password)
	if err != nil {
		log.Println(helpers.ChainError("error generating password hash", err))
		sendInternalErrorResponse(w)
		return
	}

	//save the user
	user := models.CreateNewUser(body.Username, hash)
	err = c.UserCRUD.CreateUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error saving user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}

// DeleteUser handles DELETE requests to "/user"
func (c UserController) DeleteUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	//check for id
	idStr := params.ByName("id")
	if idStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "id must be present")
		return
	}

	//parse the id
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(helpers.ChainError("error parsing user id", err))
		sendErrorResponse(w, http.StatusBadRequest, "id is in invalid format")
		return
	}

	//get the user
	user, err := c.UserCRUD.GetUserByID(id)
	if err != nil {
		log.Println(helpers.ChainError("error fetching user by id", err))
		sendInternalErrorResponse(w)
		return
	}

	if user == nil {
		sendErrorResponse(w, http.StatusBadRequest, "user not found")
		return
	}

	//delete the user
	err = c.UserCRUD.DeleteUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error deleting user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}

// PatchUserPasswordBody is the struct the body of requests to PatchUserPassword should be parsed into
type PatchUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// PatchUserPassword handles PATCH requests to "/user/password"
func (c UserController) PatchUserPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PatchUserPasswordBody

	//get the session
	sID, err := getSessionFromRequest(req)
	if err != nil {
		log.Println(helpers.ChainError("error getting session id from request", err))
		sendErrorResponse(w, http.StatusUnauthorized, "session token not provided or was in invalid format")
		return
	}

	//get the user
	user, err := c.UserCRUD.GetUserBySessionID(sID)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by session id", err))
		sendInternalErrorResponse(w)
		return
	}

	//check user was found
	if user == nil {
		sendErrorResponse(w, http.StatusUnauthorized, "no user for provided session")
		return
	}

	//parse the body
	err = parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PatchUserPassword request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//validate the body fields
	if body.OldPassword == "" || body.NewPassword == "" {
		sendErrorResponse(w, http.StatusBadRequest, "old password and new password cannot be empty")
		return
	}

	//validate old password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, body.OldPassword)
	if err != nil {
		log.Println(helpers.ChainError("error comparing password hashes", err))
		sendErrorResponse(w, http.StatusBadRequest, "old password is invalid")
		return
	}

	//validate new password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(body.NewPassword)
	if verr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(helpers.ChainError("error validating password criteria", verr))
		sendErrorResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(body.NewPassword)
	if err != nil {
		log.Println(helpers.ChainError("error generating password hash", err))
		sendInternalErrorResponse(w)
		return
	}

	//update the user
	user.PasswordHash = hash
	err = c.UserCRUD.UpdateUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error updating user", err))
		sendInternalErrorResponse(w)
		return
	}

	//return success
	sendSuccessResponse(w)
}
