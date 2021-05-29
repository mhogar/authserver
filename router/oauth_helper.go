package router

import (
	"net/http"
)

func sendOAuthErrorResponse(w http.ResponseWriter, status int, err string, description string) {
	sendResponse(w, status, OAuthErrorResponse{
		Error:            err,
		ErrorDescription: description,
	})
}
