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
	router.POST("/user", handler.PostUser)
	router.DELETE("/user/:id", handler.DeleteUser)
	router.PATCH("/user/password", handler.PatchUserPassword)

	//oauth routes
	router.POST("/oauth/token", handler.PostOAuthToken)

	return router
}
