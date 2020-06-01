package dependencies

import (
	"authserver/helpers"
	"sync"
)

var passwordCriteriaValidatorOnce sync.Once
var passwordCriteriaValidator helpers.PasswordCriteriaValidator

// ResolvePasswordCriteriaValidator resolves the PasswordCriteriaValidator dependency.
// Only the first call to this function will create a new PasswordCriteriaValidator, after which it will be retrieved from the cache.
func ResolvePasswordCriteriaValidator() helpers.PasswordCriteriaValidator {
	passwordCriteriaValidatorOnce.Do(func() {
		passwordCriteriaValidator = helpers.ConfigPasswordCriteriaValidator{}
	})
	return passwordCriteriaValidator
}
