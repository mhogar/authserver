package models_test

import (
	"testing"

	"authserver/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	Client *models.Client
}

func (suite *ClientSuite) SetupTest() {
	suite.Client = models.CreateNewClient()
}

func (suite *ClientSuite) TestCreateNewClient_CreatesClientWithSuppliedFields() {
	//act
	client := models.CreateNewClient()

	//assert
	suite.Require().NotNil(client)
	suite.NotEqual(client.ID, uuid.Nil)
}

func (suite *ClientSuite) TestValidate_WithValidClient_ReturnsValid() {
	//act
	err := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientValid, err.Status)
}

func (suite *ClientSuite) TestValidate_WithNilID_ReturnsClientInvalidID() {
	//arrange
	suite.Client.ID = uuid.Nil

	//act
	err := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientInvalidID, err.Status)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, &ClientSuite{})
}
