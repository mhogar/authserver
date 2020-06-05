package router

import (
	"authserver/controllers"
	"authserver/controllers/common"

	"github.com/julienschmidt/httprouter"
)

// CreateRouter creates a new router with the endpoints and panic handler configured
func CreateRouter(handler controllers.RequestHandler) *httprouter.Router {
	router := httprouter.New()

	router.PanicHandler = common.PanicHandler

	//user routes
	router.POST("/api/user", handler.PostAPIUser)
	router.DELETE("/api/user/:id", handler.DeleteAPIUser)
	router.PATCH("/api/user/password", handler.PatchAPIUserPassword)

	//oauth routes
	router.POST("/oauth/token", handler.PostOAuthToken)

	return router
}
