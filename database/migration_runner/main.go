package main

import (
	"authserver/common"
	"authserver/config"
	"authserver/database"
	"authserver/dependencies"
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
	down := flag.Bool("down", false, "Run migrate down instead of migrate up")
	dbKey := flag.String("db", "core", "The database to run the migrations against")
	flag.Parse()

	viper.Set("db_key", *dbKey)

	migrationRunner := migrationrunner.MigrationRunner{
		MigrationRepository: dependencies.ResolveMigrationRepository(),
		MigrationCRUD:       dependencies.ResolveDatabase(),
	}

	err = Run(dependencies.ResolveDatabase(), migrationRunner, *down)
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
func Run(db database.DBConnection, migrationRunner MigrationRunner, down bool) error {
	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		return common.ChainError("could not create database connection", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		return common.ChainError("could not reach database", err)
	}

	//run the migrations
	if down {
		err = migrationRunner.MigrateDown()
	} else {
		err = migrationRunner.MigrateUp()
	}

	if err != nil {
		return common.ChainError("error running migrations", err)
	}

	return nil
}
