package main_test

import (
	requesterror "authserver/common/request_error"
	controllermocks "authserver/controllers/mocks"
	databasemocks "authserver/database/mocks"
	"authserver/models"
	admincreator "authserver/tools/admin_creator"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AdminCreatorTestSuite struct {
	suite.Suite
	DBConnectionMock databasemocks.DBConnection
	ControllersMock  controllermocks.Controllers
}

func (suite *AdminCreatorTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.ControllersMock = controllermocks.Controllers{}
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorOpeningDatabaseConnection_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	message := "OpenConnection test error"
	suite.DBConnectionMock.On("OpenConnection").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, username, password)

	//assert
	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorPingingDatabase_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)

	message := "Ping test error"
	suite.DBConnectionMock.On("Ping").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingUser_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithNoErrors_ReturnsNoErrors() {
	//arrange
	username := "username"
	password := "password"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "OpenConnection")
	suite.DBConnectionMock.AssertCalled(suite.T(), "Ping")
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", username, password)
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.NoError(err)
	suite.NotNil(user)
}

func TestAdminCreatorTestSuite(t *testing.T) {
	suite.Run(t, &AdminCreatorTestSuite{})
}
