package main

import (
	"authserver/config"
	"authserver/database"
	"authserver/dependencies"
	commonhelpers "authserver/helpers/common"
	"flag"
	"log"

	"github.com/mhogar/migrationrunner"
	"github.com/spf13/viper"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	//parse flags
	dbKey := flag.String("db", "core", "The database to run the migrations against")
	down := flag.Bool("down", false, "Run migrate down instead of migrate up")
	flag.Parse()

	migrationRunner := migrationrunner.MigrationRunner{
		MigrationRepository: dependencies.ResolveMigrationRepository(),
		MigrationCRUD:       dependencies.ResolveDatabase(),
	}

	err = Run(dependencies.ResolveDatabase(), migrationRunner, *dbKey, *down)
	if err != nil {
		log.Fatal(err)
	}
}

// MigrationRunner is an interface to match the signature of migrationrunner's MigrationRunner.
type MigrationRunner interface {
	MigrateUp() error
	MigrateDown() error
}

// Run connects to the database and runs the migration runner. Returns any errors.
func Run(db database.DBConnection, migrationRunner MigrationRunner, dbKey string, down bool) error {
	viper.Set("db_key", dbKey)

	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		return commonhelpers.ChainError("could not create database connection", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		return commonhelpers.ChainError("could not reach database", err)
	}

	//run the migrations
	if down {
		err = migrationRunner.MigrateDown()
	} else {
		err = migrationRunner.MigrateUp()
	}

	if err != nil {
		return commonhelpers.ChainError("error running migrations", err)
	}

	return nil
}
