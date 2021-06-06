package main

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/config"
	"authserver/controllers"
	"authserver/database"
	"authserver/dependencies"
	"authserver/models"
	"flag"
	"fmt"
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

	user, err := Run(dependencies.ResolveDatabase(), dependencies.ResolveControllers(), dependencies.ResolveTransactionFactory(), *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created user:", user.ID.String())
}

// Run connects to the database and runs the admin creator. Returns any errors.
func Run(db database.DBConnection, c controllers.UserController, tf database.TransactionFactory, username string, password string) (*models.User, error) {
	//open the db connection
	err := db.OpenConnection()
	if err != nil {
		return nil, common.ChainError("could not open database connection", err)
	}

	defer db.CloseConnection()

	//check db is connected
	err = db.Ping()
	if err != nil {
		return nil, common.ChainError("could not reach database", err)
	}

	//create a new transaction
	tx, err := tf.CreateTransaction()
	if err != nil {
		return nil, err
	}

	//save the user, rollback transaction on error
	user, rerr := c.CreateUser(tx, username, password)
	if rerr.Type != requesterror.ErrorTypeNone {
		tx.RollbackTransaction()
		return nil, rerr
	}

	//commit the transaction
	err = tx.CommitTransaction()
	if err != nil {
		return nil, err
	}

	return user, nil
}
