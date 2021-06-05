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
	DBConnectionMock       databasemocks.DBConnection
	ControllersMock        controllermocks.Controllers
	TransactionFactoryMock databasemocks.TransactionFactory
	TransactionMock        databasemocks.Transaction
}

func (suite *AdminCreatorTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.ControllersMock = controllermocks.Controllers{}
	suite.TransactionFactoryMock = databasemocks.TransactionFactory{}
	suite.TransactionMock = databasemocks.Transaction{}

	suite.TransactionMock.On("RollbackTransaction")
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorOpeningDatabaseConnection_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	message := "OpenConnection test error"
	suite.DBConnectionMock.On("OpenConnection").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

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
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCreatingTransaction_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)

	message := "create transaction error"
	suite.TransactionFactoryMock.On("CreateTransaction").Return(nil, errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

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
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)

	message := "create user error"
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, requesterror.ClientError(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")
	suite.TransactionMock.AssertCalled(suite.T(), "RollbackTransaction")

	suite.Nil(user)
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *AdminCreatorTestSuite) TestRun_WithErrorCommitingTransaction_ReturnsError() {
	//arrange
	username := "username"
	password := "password"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())

	message := "commit transaction error"
	suite.TransactionMock.On("CommitTransaction").Return(errors.New(message))

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

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
	suite.TransactionFactoryMock.On("CreateTransaction").Return(&suite.TransactionMock, nil)
	suite.ControllersMock.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, requesterror.NoError())
	suite.TransactionMock.On("CommitTransaction").Return(nil)

	//act
	user, err := admincreator.Run(&suite.DBConnectionMock, &suite.ControllersMock, &suite.TransactionFactoryMock, username, password)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "OpenConnection")
	suite.DBConnectionMock.AssertCalled(suite.T(), "Ping")
	suite.TransactionFactoryMock.AssertCalled(suite.T(), "CreateTransaction")
	suite.ControllersMock.AssertCalled(suite.T(), "CreateUser", mock.Anything, username, password)
	suite.TransactionMock.AssertCalled(suite.T(), "CommitTransaction")
	suite.TransactionMock.AssertNotCalled(suite.T(), "RollbackTransaction")
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.NoError(err)
	suite.NotNil(user)
}

func TestAdminCreatorTestSuite(t *testing.T) {
	suite.Run(t, &AdminCreatorTestSuite{})
}
