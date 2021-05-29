package router

import (
	"github.com/julienschmidt/httprouter"
)

// CreateRouter creates a new router with the endpoints and panic handler configured
func CreateRouter(handlers Handlers) *httprouter.Router {
	router := httprouter.New()

	router.PanicHandler = panicHandler

	//user routes
	router.POST("/user", handlers.PostUser)
	router.DELETE("/user/:id", handlers.DeleteUser)
	router.PATCH("/user/password", handlers.PatchUserPassword)

	//token routes
	router.POST("/token", handlers.PostToken)
	router.DELETE("/token", handlers.DeleteToken)

	return router
}
