package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ScopeSuite struct {
	suite.Suite
	Scope *models.Scope
}

func (suite *ScopeSuite) SetupTest() {
	suite.Scope = models.CreateNewScope("name")
}

func (suite *ScopeSuite) TestCreateNewScope_CreatesScopeWithSuppliedFields() {
	//arrange
	name := "name"

	//act
	scope := models.CreateNewScope(name)

	//assert
	suite.Require().NotNil(scope)
	suite.NotEqual(scope.ID, uuid.Nil)
	suite.Equal(name, scope.Name)
}

func (suite *ScopeSuite) TestValidate_WithValidScope_ReturnsValid() {
	//act
	err := suite.Scope.Validate()

	//assert
	suite.EqualValues(models.ValidateScopeValid, err.Status)
}

func (suite *ScopeSuite) TestValidate_WithNilID_ReturnsScopeInvalidID() {
	//arrange
	suite.Scope.ID = uuid.Nil

	//act
	err := suite.Scope.Validate()

	//assert
	suite.EqualValues(models.ValidateScopeInvalidID, err.Status)
}

func (suite *ScopeSuite) TestValidate_WithEmptyName_ReturnsScopeInvalidName() {
	//arrange
	suite.Scope.Name = ""

	//act
	err := suite.Scope.Validate()

	//assert
	suite.EqualValues(models.ValidateScopeInvalidName, err.Status)
}

func TestScopeSuite(t *testing.T) {
	suite.Run(t, &ScopeSuite{})
}
