package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MigrationTestSuite struct {
	suite.Suite
	Migration *models.Migration
}

func (suite *MigrationTestSuite) SetupTest() {
	suite.Migration = models.CreateNewMigration(
		"00010101000000",
	)
}

func (suite *MigrationTestSuite) TestCreateNewMigration_CreatesMigrationWithSuppliedFields() {
	//arrange
	timestamp := "this is a timestamp"

	//act
	migration := models.CreateNewMigration(timestamp)

	//assert
	suite.Require().NotNil(migration)
	suite.NotEqual(migration.ID, uuid.Nil)
	suite.EqualValues(timestamp, migration.Timestamp)
}

func (suite *MigrationTestSuite) TestValidate_WithValidMigration_ReturnsModelValid() {
	//act
	err := suite.Migration.Validate()

	//assert
	suite.EqualValues(models.ValidateMigrationValid, err.Status)
}

func (suite *MigrationTestSuite) TestValidate_WithNilMigrationID_ReturnsMigrationInvalidID() {
	//arrange
	suite.Migration.ID = uuid.Nil

	//act
	err := suite.Migration.Validate()

	//assert
	suite.EqualValues(models.ValidateMigrationInvalidID, err.Status)
}

func (suite *MigrationTestSuite) TestValidate_WithVariousInvalidTimestamps_ReturnsError() {
	var timestamp string
	testCase := func() {
		//arrange
		suite.Migration.Timestamp = timestamp

		//act
		err := suite.Migration.Validate()

		//assert
		suite.EqualValues(models.ValidateMigrationInvalidTimestamp, err.Status)
	}

	timestamp = "0001010100000"
	suite.Run("TooFewDigits", testCase)

	timestamp = "000101010000000"
	suite.Run("TooManyDigits", testCase)

	timestamp = "000101010a0000"
	suite.Run("ContainsNonDigit", testCase)
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, &MigrationTestSuite{})
}
