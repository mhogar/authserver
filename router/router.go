package router

import (
	"authserver/controllers"

	"github.com/julienschmidt/httprouter"
)

type routeHandler struct {
	Control       controllers.Controllers
	Authenticator Authenticator
}

// CreateRouter creates a new router with the endpoints and panic handler configured
func CreateRouter(control controllers.Controllers, authenticator Authenticator) *httprouter.Router {
	router := httprouter.New()
	handler := routeHandler{
		Control:       control,
		Authenticator: authenticator,
	}

	router.PanicHandler = panicHandler

	//user routes
	router.POST("/user", handler.PostUser)
	router.DELETE("/user", handler.DeleteUser)
	router.PATCH("/user/password", handler.PatchUserPassword)

	//token routes
	router.POST("/token", handler.PostToken)
	router.DELETE("/token", handler.DeleteToken)

	return router
}
