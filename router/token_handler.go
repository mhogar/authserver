package router

import (
	requesterror "authserver/common/request_error"
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// PostTokenBody is the struct the body of requests to PostToken should be parsed into
type PostTokenBody struct {
	GrantType string `json:"grant_type"`
	PostTokenPasswordGrantBody
}

// PostTokenPasswordGrantBody is the struct the body of password grant requests to PostToken should be parsed into
type PostTokenPasswordGrantBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

// PostToken handles POST requests to "/token"
func (h RouteHandler) PostToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body PostTokenBody

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing PostToken request body", err))
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "invalid json body")
		return
	}

	//validate grant type is present
	if body.GrantType == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing grant_type parameter")
		return
	}

	var token *models.AccessToken = nil

	//choose the workflow based on the grant type
	switch body.GrantType {
	case "password":
		token = h.handlePasswordGrant(w, body.PostTokenPasswordGrantBody)
	default:
		sendOAuthErrorResponse(w, http.StatusBadRequest, "unsupported_grant_type", "")
	}

	if token == nil {
		return
	}

	//construct and send the access token response
	sendResponse(w, http.StatusOK, AccessTokenResponse{
		AccessToken: token.ID.String(),
		TokenType:   "bearer",
	})
}

func (h RouteHandler) handlePasswordGrant(w http.ResponseWriter, body PostTokenPasswordGrantBody) *models.AccessToken {
	//validate parameters
	if body.Username == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing username parameter")
		return nil
	}

	if body.Password == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing password parameter")
		return nil
	}

	if body.ClientID == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing client_id parameter")
		return nil
	}

	if body.Scope == "" {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_request", "missing scope parameter")
		return nil
	}

	//parse the client id
	clientID, err := uuid.Parse(body.ClientID)
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing client id", err))
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_client", "client_id was in invalid format")
		return nil
	}

	//create the token
	token, rerr := h.Control.CreateTokenFromPassword(body.Username, body.Password, clientID, body.Scope)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendOAuthErrorResponse(w, http.StatusBadRequest, rerr.ErrorName, rerr.Error())
		return nil
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return nil
	}

	return token
}

// DeleteToken handles DELETE requests to "/token"
func (h RouteHandler) DeleteToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//authenticate the user
	token, rerr := h.Authenticator.Authenticate(req)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//delete the token
	rerr = h.Control.DeleteToken(token)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusBadRequest, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//return success
	sendSuccessResponse(w)
}
