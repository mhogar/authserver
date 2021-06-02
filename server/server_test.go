package server_test

import (
	controllermocks "authserver/controllers/mocks"
	databasemocks "authserver/database/mocks"
	routermocks "authserver/router/mocks"
	"authserver/server"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	DBConnectionMock  databasemocks.DBConnection
	ControllersMock   controllermocks.Controllers
	AuthenticatorMock routermocks.Authenticator
}

func (suite *ServerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.ControllersMock = controllermocks.Controllers{}
	suite.AuthenticatorMock = routermocks.Authenticator{}
}

func (suite *ServerTestSuite) TestCreateHTTPServerRunner_CreatesRunnerUsingHTTPServer() {
	//act
	runner := server.CreateHTTPServerRunner(&suite.DBConnectionMock, &suite.ControllersMock, &suite.AuthenticatorMock)
	_, ok := runner.Server.(*server.HTTPServer)

	//assert
	suite.True(ok, "Runner's server should be an http server")
}

func (suite *ServerTestSuite) TestCreateHTTPTestServerRunner_CreatesRunnerUsingHTTPTestServer() {
	//act
	runner := server.CreateHTTPTestServerRunner(&suite.DBConnectionMock, &suite.ControllersMock, &suite.AuthenticatorMock)
	_, ok := runner.Server.(*server.HTTPTestServer)

	//assert
	suite.True(ok, "Runner's server should be an httptest server")
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}
