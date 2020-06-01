package config

import "github.com/spf13/viper"

// DatabaseConfig is a struct with fields needed for configuring database operations.
type DatabaseConfig struct {
	// URL is the url of the database server.
	URL string

	// Port is the port on the database server to connect to.
	Port string

	// Timeout is the default timeout all database requests should use.
	Timeout int

	// DBs is a string map that maps db keys to the actual name of the database.
	DBs map[string]string
}

func initDatabaseConfig() {
	config := make(map[string]interface{})

	config["local"] = DatabaseConfig{
		URL:     "localhost",
		Port:    "27017",
		Timeout: 3000,
		DBs: map[string]string{
			"core":        "core",
			"integration": "integration",
		},
	}

	config["travis"] = DatabaseConfig{
		URL:     "127.0.0.1",
		Port:    "27017",
		Timeout: 3000,
		DBs: map[string]string{
			"core":        "core",
			"integration": "integration",
		},
	}

	viper.Set("database", config)
}
