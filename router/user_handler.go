package router

import (
	"log"
	"net/http"

	"authserver/common"
	requesterror "authserver/common/request_error"

	"github.com/julienschmidt/httprouter"
)

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (h routeHandler) PostUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//parse the body
	var body PostUserBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostUser request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//create the user
	_, rerr := h.Control.CreateUser(body.Username, body.Password)
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

// DeleteUser handles DELETE requests to "/user"
func (h routeHandler) DeleteUser(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//authenticate the user
	token, rerr := h.Authenticator.Authenticate(req)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//delete the user
	rerr = h.Control.DeleteUser(token.User)
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

// PatchUserPasswordBody is the struct the body of requests to PatchUserPassword should be parsed into
type PatchUserPasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// PatchUserPassword handles PATCH requests to "/user/password"
func (h routeHandler) PatchUserPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//authenticate the user
	token, rerr := h.Authenticator.Authenticate(req)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//parse the body
	var body PatchUserPasswordBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PatchUserPassword request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//update the password
	rerr = h.Control.UpdateUserPassword(token.User, body.OldPassword, body.NewPassword)
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
