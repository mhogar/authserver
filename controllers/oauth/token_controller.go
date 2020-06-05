package oauth

import (
	"authserver/controllers/common"
	"authserver/database"
	"authserver/helpers"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TokenController handles requests to "/oauth/token" endpoints
type TokenController struct {
	UserCRUD       database.UserCRUD
	SessionCRUD    database.SessionCRUD
	PasswordHasher helpers.PasswordHasher
	TokenHelper    helpers.TokenHelper
}

// PostTokenBody is the struct the body of requests to PostToken should be parsed into
type PostTokenBody struct {
	GrantType string `json:"grant_type"`
	passwordGrantBody
	authorizationCodeGrantBody
}

//PostOAuthToken handles Post requests to "/oauth/token"
func (c TokenController) PostOAuthToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostTokenBody

	//parse the body
	err := common.ParseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(helpers.ChainError("error parsing PostToken request body", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//choose the workflow based on the grant type
	switch body.GrantType {
	case "password":
		c.handlePasswordGrant(w, body.passwordGrantBody)
	default:
	}
}

func (c TokenController) handlePasswordGrant(w http.ResponseWriter, body passwordGrantBody) {
	//get the user
	user, err := c.UserCRUD.GetUserByUsername(body.Username)
	if err != nil {
		log.Println(helpers.ChainError("error getting user by username", err))
		common.SendInternalErrorResponse(w)
		return
	}

	if user == nil {
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return
	}

	//validate the password
	err = c.PasswordHasher.ComparePasswords(user.PasswordHash, body.Password)
	if err != nil {
		log.Println(helpers.ChainError("error comparing password hashes", err))
		common.SendErrorResponse(w, http.StatusBadRequest, "invalid username and/or password")
		return
	}

	//create new session
	session := models.CreateNewSession(user.ID)

	//generate token
	token, err := c.TokenHelper.CreateSessionToken(session.ID)
	if err != nil {
		log.Println(helpers.ChainError("error creating session token", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//save the session
	err = c.SessionCRUD.CreateSession(session)
	if err != nil {
		log.Println(helpers.ChainError("error saving session", err))
		common.SendInternalErrorResponse(w)
		return
	}

	//return the token
	common.SendDataResponse(w, token)
}
