package controllers

import (
	"authserver/controllers/api"
	"authserver/controllers/oauth"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RequestHandler is an interface that encapsulates all other handler interfaces
type RequestHandler interface {
	UserHandler
	TokenHandler
}

// UserHandler is an interface for handling requests to user routes
type UserHandler interface {
	// PostAPIUser handles POST requests to "/user"
	PostAPIUser(http.ResponseWriter, *http.Request, httprouter.Params)

	// DeleteAPIUser handles DELETE requests to "/user/:id"
	DeleteAPIUser(http.ResponseWriter, *http.Request, httprouter.Params)

	// PatchAPIUserPassword handles PATCH requests to "/user/password"
	PatchAPIUserPassword(http.ResponseWriter, *http.Request, httprouter.Params)
}

// TokenHandler is an interface for handling requests to token routes
type TokenHandler interface {
	//PostOAuthToken handles Post requests to "/token"
	PostOAuthToken(http.ResponseWriter, *http.Request, httprouter.Params)
}

// RequestHandle is an implementation of RequestHandler that uses controllers to satisfy the interface's methods
type RequestHandle struct {
	api.UserController
	oauth.TokenController
}
