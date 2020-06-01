package dependencies

import (
	"authserver/helpers"
	"sync"
)

var tokenHelperOnce sync.Once
var tokenHelper helpers.TokenHelper

// ResolveTokenHelper resolves the TokenHelper dependency.
// Only the first call to this function will create a new TokenHelper, after which it will be retrieved from the cache.
func ResolveTokenHelper() helpers.TokenHelper {
	tokenHelperOnce.Do(func() {
		tokenHelper = helpers.JWTHelper{}
	})
	return tokenHelper
}
