package dependencies

import (
	"authserver/router"
	"sync"
)

var createHandlersOnce sync.Once
var handlers router.Handlers

// ResolveHandlers resolves the Handlers dependency.
// Only the first call to this function will create a new Handlers, after which it will be retrieved from the cache.
func ResolveHandlers() router.Handlers {
	createHandlersOnce.Do(func() {
		handlers = &router.Handles{
			UserHandle: router.UserHandle{
				UserControl:   ResolveControllers(),
				Authenticator: ResolveAuthenticator(),
			},
			TokenHandle: router.TokenHandle{},
		}
	})
	return handlers
}
