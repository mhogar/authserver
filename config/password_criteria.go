package config

import "github.com/spf13/viper"

// PasswordCriteriaConfig is a struct for encapsulating criteria requirements for a password
type PasswordCriteriaConfig struct {
	// MinLength is the minimum length the password must be
	MinLength int

	// RequireLowerCase determines if at least one lower case letter must be present
	RequireLowerCase bool

	// RequireUpperCase determines if at least one upper case letter must be present
	RequireUpperCase bool

	// RequireDigit determines if at least one digit must be present
	RequireDigit bool

	// RequireSymbol determines if at least one symbol must be present
	RequireSymbol bool
}

func initPasswordCriteriaConfig() {
	viper.Set("password", PasswordCriteriaConfig{
		MinLength:        8,
		RequireLowerCase: true,
		RequireUpperCase: true,
		RequireDigit:     true,
		RequireSymbol:    true,
	})
}
