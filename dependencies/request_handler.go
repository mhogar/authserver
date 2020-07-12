package dependencies

import (
	"authserver/controllers"
	"sync"
)

var requestHandlerOnce sync.Once
var requestHandler controllers.RequestHandler

// ResolveRequestHandler resolves the RouteHandler dependency.
// Only the first call to this function will create a new RouteHandler, after which it will be retrieved from the cache.
func ResolveRequestHandler() controllers.RequestHandler {
	requestHandlerOnce.Do(func() {
		requestHandler = &controllers.RequestHandle{
			UserController: controllers.UserController{
				CRUD:                      ResolveDatabase(),
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			},
			TokenController: controllers.TokenController{
				CRUD:           ResolveDatabase(),
				PasswordHasher: ResolvePasswordHasher(),
			},
		}
	})
	return requestHandler
}
