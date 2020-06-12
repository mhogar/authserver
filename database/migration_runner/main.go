package main

import (
	"authserver/config"
	"authserver/database"
	"authserver/dependencies"
	"authserver/helpers"
	"flag"
	"log"

	"github.com/mhogar/migrationrunner"
)

func main() {
	config.InitConfig()

	//parse flags
	down := flag.Bool("down", false, "Run migrate down instead of migrate up")
	flag.Parse()

	err := Run(dependencies.ResolveDatabase(), dependencies.ResolveDatabase(), dependencies.ResolveMigrationRepository(), *down)
	if err != nil {
		log.Fatal(err)
	}
}

// Run connects to the database and runs the migration runner. Returns any errors.
func Run(db database.DBConnection, migrationCRUD migrationrunner.MigrationCRUD, repo migrationrunner.MigrationRepository, down bool) error {
	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		return helpers.ChainError("could not create database connection", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		return helpers.ChainError("could not reach database", err)
	}

	//run the migrations
	if down {
		err = migrationrunner.MigrateDown(repo, migrationCRUD)
	} else {
		err = migrationrunner.MigrateUp(repo, migrationCRUD)
	}

	if err != nil {
		return helpers.ChainError("error running migrations", err)
	}

	return nil
}
