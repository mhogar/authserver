package main

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/config"
	"authserver/controllers"
	"authserver/database"
	"authserver/dependencies"
	"flag"
	"log"

	"github.com/spf13/viper"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	//parse flags
	dbKey := flag.String("db", "core", "The database to run the scipt against")
	username := flag.String("username", "", "The username for the admin")
	password := flag.String("password", "", "The password for the admin")
	flag.Parse()

	viper.Set("db_key", *dbKey)

	err = Run(dependencies.ResolveDatabase(), dependencies.ResolveControllers(), *username, *password)
	if err != nil {
		log.Fatal(err)
	}
}

// Run connects to the database and runs the admin creator. Returns any errors.
func Run(db database.DBConnection, c controllers.UserController, username string, password string) error {
	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		return common.ChainError("could not open database connection", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		return common.ChainError("could not reach database", err)
	}

	//save the user
	_, rerr := c.CreateUser(username, password)
	if rerr.Type != requesterror.ErrorTypeNone {
		return rerr
	}

	return nil
}
