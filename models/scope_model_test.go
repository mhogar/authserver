package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ScopeTestSuite struct {
	suite.Suite
	Scope *models.Scope
}

func (suite *ScopeTestSuite) SetupTest() {
	suite.Scope = models.CreateNewScope("name")
}

func (suite *ScopeTestSuite) TestCreateNewScope_CreatesScopeWithSuppliedFields() {
	//arrange
	name := "name"

	//act
	scope := models.CreateNewScope(name)

	//assert
	suite.Require().NotNil(scope)
	suite.NotEqual(scope.ID, uuid.Nil)
	suite.Equal(name, scope.Name)
}

func (suite *ScopeTestSuite) TestValidate_WithValidScope_ReturnsValid() {
	//act
	err := suite.Scope.Validate()

	//assert
	suite.Equal(models.ValidateScopeValid, err.Status)
}

func (suite *ScopeTestSuite) TestValidate_WithNilID_ReturnsScopeNilID() {
	//arrange
	suite.Scope.ID = uuid.Nil

	//act
	err := suite.Scope.Validate()

	//assert
	suite.Equal(models.ValidateScopeNilID, err.Status)
}

func (suite *ScopeTestSuite) TestValidate_WithEmptyName_ReturnsScopeInvalidName() {
	//arrange
	suite.Scope.Name = ""

	//act
	err := suite.Scope.Validate()

	//assert
	suite.Equal(models.ValidateScopeEmptyName, err.Status)
}

func (suite *ScopeTestSuite) TestValidate_ScopeNameMaxLengthTestCases() {
	var name string
	var expectedValidateError int

	testCase := func() {
		//arrange
		suite.Scope.Name = name

		//act
		err := suite.Scope.Validate()

		//assert
		suite.Equal(expectedValidateError, err.Status)
	}

	name = "aaaaaaaaaaaaaaa" //15 chars
	expectedValidateError = models.ValidateScopeValid
	suite.Run("ExactlyMaxLengthIsValid", testCase)

	name = "aaaaaaaaaaaaaaaa" //16 chars
	expectedValidateError = models.ValidateScopeNameTooLong
	suite.Run("OneMoreThanMaxLengthIsInvalid", testCase)
}

func TestScopeTestSuite(t *testing.T) {
	suite.Run(t, &ScopeTestSuite{})
}
