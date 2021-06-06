package server_test

import (
	"authserver/common"
	databasemocks "authserver/database/mocks"
	"authserver/server"
	"authserver/server/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite
	DBConnectionMock databasemocks.DBConnection
	ServerMock       mocks.Server
	Runner           *server.Runner
}

func (suite *RunnerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.ServerMock = mocks.Server{}

	suite.Runner = &server.Runner{
		DBConnection: &suite.DBConnectionMock,
		Server:       &suite.ServerMock,
	}
}

func (suite *RunnerTestSuite) TestRun_WithErrorOpeningDBConnection_ReturnsError() {
	//arrange
	message := "OpenConnection mock error"
	suite.DBConnectionMock.On("OpenConnection").Return(errors.New(message))

	//act
	err := suite.Runner.Run()

	//assert
	common.AssertError(&suite.Suite, err, message)
}

func (suite *RunnerTestSuite) TestRun_WithPingingDatabase_ReturnsError() {
	//arrange
	message := "Ping mock error"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(errors.New(message))

	//act
	err := suite.Runner.Run()

	//assert
	common.AssertError(&suite.Suite, err, message)
}

func (suite *RunnerTestSuite) TestRun_WithErrorStartingServer_ReturnsError() {
	//arrange
	message := "Start mock error"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.ServerMock.On("Start").Return(errors.New(message))

	//act
	err := suite.Runner.Run()

	//assert
	common.AssertError(&suite.Suite, err, message)
}

func (suite *RunnerTestSuite) TestRun_StartsServer() {
	//arrange
	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.ServerMock.On("Start").Return(nil)

	//act
	err := suite.Runner.Run()

	//assert
	suite.Require().NoError(err)

	suite.DBConnectionMock.AssertCalled(suite.T(), "OpenConnection")
	suite.DBConnectionMock.AssertCalled(suite.T(), "Ping")
	suite.ServerMock.AssertCalled(suite.T(), "Start")
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, &RunnerTestSuite{})
}
