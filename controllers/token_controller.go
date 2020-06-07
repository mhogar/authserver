package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TokenController handles requests to "/token" endpoints
type TokenController struct{}

// PostToken handles POST requests to "/token"
func (c TokenController) PostToken(http.ResponseWriter, *http.Request, httprouter.Params) {
}
