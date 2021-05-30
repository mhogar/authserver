package main_test

import (
	databasemocks "authserver/database/mocks"
	migrationrunner "authserver/tools/migration_runner"
	"authserver/tools/migration_runner/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MigrationRunnerTestSuite struct {
	suite.Suite
	DBConnectionMock    databasemocks.DBConnection
	MigrationRunnerMock mocks.MigrationRunner
}

func (suite *MigrationRunnerTestSuite) SetupTest() {
	suite.DBConnectionMock = databasemocks.DBConnection{}
	suite.MigrationRunnerMock = mocks.MigrationRunner{}
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorOpeningDatabaseConnection_ReturnsError() {
	//arrange
	message := "OpenConnection test error"

	suite.DBConnectionMock.On("OpenConnection").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, false)

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
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, false)

	//assert
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownFalse_RunsUpMigration() {
	//arrange
	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.MigrationRunnerMock.On("MigrateUp").Return(nil)

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, false)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.NoError(err)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorRunningUpMigration_ReturnsError() {
	//arrange
	message := "MigrateUp test error"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.MigrationRunnerMock.On("MigrateUp").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, false)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateUp")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateDown")
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithDownTrue_RunsDownMigration() {
	//arrange
	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.MigrationRunnerMock.On("MigrateDown").Return(nil)

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, true)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.NoError(err)
}

func (suite *MigrationRunnerTestSuite) TestRun_WithErrorRunningDownMigration_ReturnsError() {
	//arrange
	message := "MigrateDown test error"

	suite.DBConnectionMock.On("OpenConnection").Return(nil)
	suite.DBConnectionMock.On("CloseConnection").Return(nil)
	suite.DBConnectionMock.On("Ping").Return(nil)
	suite.MigrationRunnerMock.On("MigrateDown").Return(errors.New(message))

	//act
	err := migrationrunner.Run(&suite.DBConnectionMock, &suite.MigrationRunnerMock, true)

	//assert
	suite.MigrationRunnerMock.AssertCalled(suite.T(), "MigrateDown")
	suite.MigrationRunnerMock.AssertNotCalled(suite.T(), "MigrateUp")
	suite.DBConnectionMock.AssertCalled(suite.T(), "CloseConnection")

	suite.Require().Error(err)
	suite.Contains(err.Error(), message)
}

func TestMigrationRunnerTestSuite(t *testing.T) {
	suite.Run(t, &MigrationRunnerTestSuite{})
}
