package router

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	requesterror "authserver/common/request_error"
	"authserver/controllers"
	commonhelpers "authserver/helpers/common"

	"github.com/julienschmidt/httprouter"
)

// UserControl handles requests to "/user" endpoints
type UserHandle struct {
	UserControl   controllers.UserController
	Authenticator Authenticator
}

// PostUserBody is the struct the body of requests to PostUser should be parsed into
type PostUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostUser handles Post requests to "/user"
func (h UserHandle) PostUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//authenticate the user
	_, rerr := h.Authenticator.Authenticate(req)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//parse the body
	var body PostUserBody
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing PostUser request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//create the user
	_, rerr = h.UserControl.CreateUser(body.Username, body.Password)
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
func (h UserHandle) DeleteUser(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//authenticate the user
	_, rerr := h.Authenticator.Authenticate(req)
	if rerr.Type == requesterror.ErrorTypeClient {
		sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
		return
	} else if rerr.Type == requesterror.ErrorTypeInternal {
		sendInternalErrorResponse(w, rerr.Error())
		return
	}

	//check for id
	idStr := params.ByName("id")
	if idStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "id must be present")
		return
	}

	//parse the id
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing user id", err))
		sendErrorResponse(w, http.StatusBadRequest, "id is in invalid format")
		return
	}

	//delete the user
	rerr = h.UserControl.DeleteUser(id)
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
func (h UserHandle) PatchUserPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
		log.Println(commonhelpers.ChainError("error parsing PatchUserPassword request body", err))
		sendErrorResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	//update the password
	rerr = h.UserControl.UpdateUserPassword(token.User, body.OldPassword, body.NewPassword)
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
