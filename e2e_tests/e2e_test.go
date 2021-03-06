package e2e_test

import (
	"authserver/common"
	"authserver/config"
	"authserver/database"
	"authserver/dependencies"
	"authserver/server"
	"net/http"
	"net/http/httptest"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	DBConnection database.DBConnection
	Server       *httptest.Server
}

func (suite *E2ETestSuite) SetupSuite() {
	err := config.InitConfig("../")
	suite.Require().NoError(err)

	//set db key and create database
	viper.Set("db_key", "integration")
	suite.DBConnection = dependencies.ResolveDatabase()

	//create the test server
	runner := server.CreateHTTPTestServerRunner(suite.DBConnection, dependencies.ResolveRouterFactory())
	suite.Server = runner.Server.(*server.HTTPTestServer).Server

	// run the server
	err = runner.Run()
	suite.Require().NoError(err)
}

func (suite *E2ETestSuite) SendRequest(method string, endpoint string, bearerToken string, body interface{}) *http.Response {
	req := common.CreateRequest(&suite.Suite, method, suite.Server.URL+endpoint, bearerToken, body)

	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	return res
}
