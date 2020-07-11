package server_test

import (
	controllermocks "authserver/controllers/mocks"
	databasemocks "authserver/database/mocks"
	"authserver/server"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	DBConnectionMock   databasemocks.DBConnection
	RequestHandlerMock controllermocks.RequestHandler
}

func (suite *ServerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.RequestHandlerMock = controllermocks.RequestHandler{}
}

func (suite *ServerTestSuite) TestCreateHTTPServerRunner_CreatesRunnerUsingHTTPServer() {
	//act
	runner := server.CreateHTTPServerRunner(&suite.DBConnectionMock, &suite.RequestHandlerMock)
	_, ok := runner.Server.(*server.HTTPServer)

	//assert
	suite.True(ok, "Runner's server should be an http server")
}

func (suite *ServerTestSuite) TestCreateHTTPTestServerRunner_CreatesRunnerUsingHTTPTestServer() {
	//act
	runner := server.CreateHTTPTestServerRunner(&suite.DBConnectionMock, &suite.RequestHandlerMock)
	_, ok := runner.Server.(*server.HTTPTestServer)

	//assert
	suite.True(ok, "Runner's server should be an httptest server")
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}
