package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handlers interface {
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

	// DeleteToken handles DELETE requests to "/token"
	DeleteToken(http.ResponseWriter, *http.Request, httprouter.Params)
}

type Handles struct {
	UserHandle
	TokenHandle
}
