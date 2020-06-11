package config

import "github.com/spf13/viper"

// DatabaseConfig is a struct with fields needed for configuring database operations.
type DatabaseConfig struct {
	// DBs is a string map that maps db keys to the connection string of the database.
	DBs map[string]string

	// Timeout is the default timeout all database requests should use.
	Timeout int
}

func initDatabaseConfig() {
	config := make(map[string]interface{})

	config["local"] = DatabaseConfig{
		DBs: map[string]string{
			"core":        "postgres://postgres:password@localhost/postgres?sslmode=disable",
			"integration": "postgres://postgres:password@localhost/postgres?sslmode=disable",
		},
		Timeout: 3000,
	}

	config["travis"] = DatabaseConfig{
		DBs: map[string]string{
			"integration": "postgres://postgres:@localhost/travis_ci_test?sslmode=disable",
		},
		Timeout: 3000,
	}

	viper.Set("database", config)
}
