package main_test

import (
	migrationrunner "authserver/database/migration_runner"
	databasemocks "authserver/database/mocks"
	"errors"
	"testing"

	migrationrunnermocks "github.com/mhogar/migrationrunner/mocks"
	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	DBConnectionMock        databasemocks.DBConnection
	MigrationCRUDMock       migrationrunnermocks.MigrationCRUD
	MigrationRepositoryMock migrationrunnermocks.MigrationRepository
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.MigrationCRUDMock = migrationrunnermocks.MigrationCRUD{}
	suite.MigrationRepositoryMock = migrationrunnermocks.MigrationRepository{}
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorOpeningDatabaseConnection_ReturnsError() {
	//arrange
	message := "OpenConnection test error"

	suite.DBConnectionMock.On("OpenConnection").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationCRUDMock, &suite.MigrationRepositoryMock, false)

	//assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorPingingDatabase_ReturnsError() {
	//arrange
	message := "Ping test error"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationCRUDMock, &suite.MigrationRepositoryMock, false)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}
