package config

import "github.com/spf13/viper"

//InitConfig sets the default config values and binds environment variables. Should be called at the start of the application.
func InitConfig() {
	viper.SetDefault("env", "local")

	viper.SetEnvPrefix("cfg")
	viper.BindEnv("env")

	initDatabaseConfig()
	initPasswordCriteriaConfig()
}
