package server_test

import (
	databasemocks "authserver/database/mocks"
	routermocks "authserver/router/mocks"
	"authserver/server"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	DBConnectionMock  databasemocks.DBConnection
	RouterFactoryMock routermocks.IRouterFactory
}

func (suite *ServerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.RouterFactoryMock = routermocks.IRouterFactory{}
}

func (suite *ServerTestSuite) TestCreateHTTPServerRunner_CreatesRunnerUsingHTTPServer() {
	//arrange
	suite.RouterFactoryMock.On("CreateRouter").Return(nil)

	//act
	runner := server.CreateHTTPServerRunner(&suite.DBConnectionMock, &suite.RouterFactoryMock)
	_, ok := runner.Server.(*server.HTTPServer)

	//assert
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
	suite.True(ok, "Runner's server should be an http server")
}

func (suite *ServerTestSuite) TestCreateHTTPTestServerRunner_CreatesRunnerUsingHTTPTestServer() {
	//arrange
	suite.RouterFactoryMock.On("CreateRouter").Return(nil)

	//act
	runner := server.CreateHTTPTestServerRunner(&suite.DBConnectionMock, &suite.RouterFactoryMock)
	_, ok := runner.Server.(*server.HTTPTestServer)

	//assert
	suite.RouterFactoryMock.AssertCalled(suite.T(), "CreateRouter")
	suite.True(ok, "Runner's server should be an httptest server")
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &ServerTestSuite{})
}
