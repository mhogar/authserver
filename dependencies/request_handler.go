package dependencies

import (
	"authserver/controllers"
	"authserver/controllers/api"
	"authserver/controllers/oauth"
	"sync"
)

var requestHandlerOnce sync.Once
var requestHandler controllers.RequestHandler

// ResolveRequestHandler resolves the RouteHandler dependency.
// Only the first call to this function will create a new RouteHandler, after which it will be retrieved from the cache.
func ResolveRequestHandler() controllers.RequestHandler {
	requestHandlerOnce.Do(func() {
		requestHandler = &controllers.RequestHandle{
			UserController: api.UserController{
				UserCRUD:                  ResolveDatabase(),
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			},
			TokenController: oauth.TokenController{
				UserCRUD:       ResolveDatabase(),
				SessionCRUD:    ResolveDatabase(),
				PasswordHasher: ResolvePasswordHasher(),
				TokenHelper:    ResolveTokenHelper(),
			},
		}
	})
	return requestHandler
}
