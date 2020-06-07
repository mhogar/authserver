package controllers

import (
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
	// PostUser handles POST requests to "/user"
	PostUser(http.ResponseWriter, *http.Request, httprouter.Params)

	// DeleteUser handles DELETE requests to "/user/:id"
	DeleteUser(http.ResponseWriter, *http.Request, httprouter.Params)

	// PatchUserPassword handles PATCH requests to "/user/password"
	PatchUserPassword(http.ResponseWriter, *http.Request, httprouter.Params)
}

// TokenHandler is an interface for handling requests to token routes
type TokenHandler interface {
	// PostToken handles POST requests to "/token"
	PostToken(http.ResponseWriter, *http.Request, httprouter.Params)
}

// RequestHandle is an implementation of RequestHandler that uses controllers to satisfy the interface's methods
type RequestHandle struct {
	UserController
	TokenController
}
