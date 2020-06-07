package controllers

import (
	"log"
	"net/http"
)

// BasicResponse represents a response with a simple true/false success field
type BasicResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse represents a response with a true/false success field and an error message
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// DataResponse represents a response with a true/false success field and generic data
type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// OAuthErrorResponse represents an error response defined by the oauth spec
type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// AccessTokenResponse represents an access token response defined by the oauth spec
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// PanicHandler is the function to be called if a panic is encountered
func PanicHandler(w http.ResponseWriter, req *http.Request, info interface{}) {
	log.Println(info)

	sendInternalErrorResponse(w)
}
