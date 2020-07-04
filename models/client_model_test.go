package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	Client *models.Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.Client = models.CreateNewClient()
}

func (suite *ClientTestSuite) TestCreateNewClient_CreatesClientWithSuppliedFields() {
	//act
	client := models.CreateNewClient()

	//assert
	suite.Require().NotNil(client)
	suite.NotEqual(client.ID, uuid.Nil)
}

func (suite *ClientTestSuite) TestValidate_WithValidClient_ReturnsValid() {
	//act
	err := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientValid, err.Status)
}

func (suite *ClientTestSuite) TestValidate_WithNilID_ReturnsClientInvalidID() {
	//arrange
	suite.Client.ID = uuid.Nil

	//act
	err := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientInvalidID, err.Status)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
