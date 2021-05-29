package main

import (
	"authserver/dependencies"
	commonhelpers "authserver/helpers/common"
	"authserver/server"
	"log"

	"authserver/config"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(commonhelpers.ChainError("error initing config", err))
	}

	serverRunner := server.CreateHTTPServerRunner(dependencies.ResolveDatabase(), dependencies.ResolveHandlers())
	log.Fatal(serverRunner.Run())
}
