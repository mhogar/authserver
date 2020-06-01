package main

import (
	"authserver/dependencies"
	"authserver/router"
	"fmt"
	"log"
	"net/http"

	"authserver/config"
)

func main() {
	config.InitConfig()

	//create the server
	server := &http.Server{
		Addr:    ":8443",
		Handler: router.CreateRouter(dependencies.ResolveRequestHandler()),
	}

	fmt.Println("Server is running on port", server.Addr)

	//run the server
	log.Fatal(server.ListenAndServe())
}
