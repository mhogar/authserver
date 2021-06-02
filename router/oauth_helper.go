package router

import (
	"authserver/common"
	"net/http"
)

func sendOAuthErrorResponse(w http.ResponseWriter, status int, err string, description string) {
	sendResponse(w, status, common.OAuthErrorResponse{
		Error:            err,
		ErrorDescription: description,
	})
}
