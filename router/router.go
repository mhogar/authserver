package router

import (
	"authserver/controllers"
	"authserver/database"

	"github.com/julienschmidt/httprouter"
)

type IRouterFactory interface {
	CreateRouter() *httprouter.Router
}

type RouterFactory struct {
	Controllers        controllers.Controllers
	Authenticator      Authenticator
	TransactionFactory database.TransactionFactory
}

// CreateRouter creates a new httprouter with the endpoints and panic handler configured.
func (rf RouterFactory) CreateRouter() *httprouter.Router {
	r := httprouter.New()
	r.PanicHandler = panicHandler

	//user routes
	r.POST("/user", rf.createHandler(rf.postUser, false))
	r.DELETE("/user", rf.createHandler(rf.deleteUser, true))
	r.PATCH("/user/password", rf.createHandler(rf.patchUserPassword, true))

	//token routes
	r.POST("/token", rf.createHandler(rf.postToken, false))
	r.DELETE("/token", rf.createHandler(rf.deleteToken, true))

	return r
}
