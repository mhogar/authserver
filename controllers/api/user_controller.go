package api

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"authserver/controllers/common"
	"authserver/database"
	"authserver/helpers"
	"authserver/models"

	"github.com/julienschmidt/httprouter"
)

// UserController handles requests to "/api/user" endpoints
type UserController struct {
	UserCRUD                  database.UserCRUD
	PasswordHasher            helpers.PasswordHasher
	PasswordCriteriaValidator helpers.PasswordCriteriaValidator
}

// PostAPIUserBody is the struct the body of requests to PostAPIUser should be parsed into
type PostAPIUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostAPIUser handles Post requests to "/api/user"
func (c *UserController) PostAPIUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostAPIUserBody

	//parse the body
	err := common.ParseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PostAPIUser request body", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//validate the body fields
	if body.Username == "" || body.Password == "" {
		common.SendErrorResponse(w, http.StatusBadRequest, "username and password cannot be empty")
		return
	}

	//validate username is unique
	otherUser, err := c.UserCRUD.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by username", err))
		common.SendInternalErrorResponse(w)
		return
	}

	if otherUser != nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "an user with that username already exists")
		return
	}

	//validate password meets criteria
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(body.Password)
	if verr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(helpers.ChainError("error validating password criteria", verr))
		common.SendErrorResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(body.Password)
	if err != nil {
		log.Println(helpers.ChainError("error generating password hash", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//save the user
	user := models.CreateNewUser(body.Username, hash)
	err = c.UserCRUD.CreateUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error saving user", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//return success
	common.SendSuccessResponse(w)
}

// DeleteAPIUser handles DELETE requests to "/api/user"
func (c *UserController) DeleteAPIUser(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	//check for id
	idStr := params.ByName("id")
	if idStr == "" {
		common.SendErrorResponse(w, http.StatusBadRequest, "id must be present")
		return
	}

	//parse the id
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(helpers.ChainError("error parsing user id", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "id is in invalid format")
		return
	}

	//get the user
	user, err := c.UserCRUD.GetUserByID(id)
	if err != nil {
		log.Println(helpers.ChainError("error fetching user by id", err))
		common.SendInternalErrorResponse(w)
		return
	}

	if user == nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "user not found")
		return
	}

	//delete the user
	err = c.UserCRUD.DeleteUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error deleting user", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//return success
	common.SendSuccessResponse(w)
}

// PatchAPIUserPasswordBody is the struct the body of requests to PatchAPIUserPassword should be parsed into
type PatchAPIUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// PatchAPIUserPassword handles PATCH requests to "/api/user/password"
func (c *UserController) PatchAPIUserPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PatchAPIUserPasswordBody

	//get the session
	sID, err := getSessionFromRequest(req)
	if err != nil {
		log.Println(helpers.ChainError("error getting session id from request", err))
		common.SendErrorResponse(w, http.StatusUnauthorized, "session token not provided or was in invalid format")
		return
	}

	//get the user
	user, err := c.UserCRUD.GetUserBySessionID(sID)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by session id", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//check user was found
	if user == nil {
		common.SendErrorResponse(w, http.StatusUnauthorized, "no user for provided session")
		return
	}

	//parse the body
	err = common.ParseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PatchAPIUserPassword request body", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//validate the body fields
	if body.OldPassword == "" || body.NewPassword == "" {
		common.SendErrorResponse(w, http.StatusBadRequest, "old password and new password cannot be empty")
		return
	}

	//validate old password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, body.OldPassword)
	if err != nil {
		log.Println(helpers.ChainError("error comparing password hashes", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "old password is invalid")
		return
	}

	//validate new password meets critera
	verr := c.PasswordCriteriaValidator.ValidatePasswordCriteria(body.NewPassword)
	if verr.Status != helpers.ValidatePasswordCriteriaValid {
		log.Println(helpers.ChainError("error validating password criteria", verr))
		common.SendErrorResponse(w, http.StatusBadRequest, "password does not meet minimum criteria")
		return
	}

	//hash the password
	hash, err := c.PasswordHasher.HashPassword(body.NewPassword)
	if err != nil {
		log.Println(helpers.ChainError("error generating password hash", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//update the user
	user.PasswordHash = hash
	err = c.UserCRUD.UpdateUser(user)
	if err != nil {
		log.Println(helpers.ChainError("error updating user", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//return success
	common.SendSuccessResponse(w)
}
