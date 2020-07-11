package main

import (
	"authserver/dependencies"
	commonhelpers "authserver/helpers/common"
	"authserver/router"
	"fmt"
	"log"
	"net/http"

	"authserver/config"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(commonhelpers.ChainError("error initing config", err))
	}

	//connect to the database
	db := dependencies.ResolveDatabase()

	err = db.OpenConnection()
	if err != nil {
		log.Fatal(commonhelpers.ChainError("error opening database connection", err))
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(commonhelpers.ChainError("error reaching database", err))
	}

	//create the server
	server := &http.Server{
		Addr:    ":8443",
		Handler: router.CreateRouter(dependencies.ResolveRequestHandler()),
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServe())
}
