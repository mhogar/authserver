package controllers

import (
	"authserver/database"
	"authserver/helpers"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TokenController handles requests to "/token" endpoints
type TokenController struct {
	UserCRUD       database.UserCRUD
	SessionCRUD    database.SessionCRUD
	PasswordHasher helpers.PasswordHasher
	TokenHelper    helpers.TokenHelper
}

//PostToken handles Post requests to "/token"
func (c TokenController) PostToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PostToken request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//get the user
	user, err := c.UserCRUD.GetUserByUsername(body.Email)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by username", err))
		sendInternalErrorResponse(w)
		return
	}

	if user == nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, body.Password)
	if err != nil {
		log.Println(helpers.ChainError("error comparing password hashes", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return
	}

	//create new session
	session := models.CreateNewSession(user.ID)

	//generate token
	token, err := c.TokenHelper.CreateSessionToken(session.ID)
	if err != nil {
		log.Println(helpers.ChainError("error creating session token", err))
		sendInternalErrorResponse(w)
		return
	}

	//save the session
	err = c.SessionCRUD.CreateSession(session)
	if err != nil {
		log.Println(helpers.ChainError("error saving session", err))
		sendInternalErrorResponse(w)
		return
	}

	//return the token
	sendDataResponse(w, token)
}
