package dependencies

import (
	"authserver/helpers"
	"sync"
)

var passwordHasherOnce sync.Once
var passwordHasher helpers.PasswordHasher

// ResolvePasswordHasher resolves the PasswordHasher dependency.
// Only the first call to this function will create a new PasswordHasher, after which it will be retrieved from the cache.
func ResolvePasswordHasher() helpers.PasswordHasher {
	passwordHasherOnce.Do(func() {
		passwordHasher = helpers.BCryptPasswordHasher{}
	})
	return passwordHasher
}
