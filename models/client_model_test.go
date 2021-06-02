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
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientValid, verr)
}

func (suite *ClientTestSuite) TestValidate_WithNilID_ReturnsClientNilID() {
	//arrange
	suite.Client.ID = uuid.Nil

	//act
	verr := suite.Client.Validate()

	//assert
	suite.Equal(models.ValidateClientNilID, verr)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}
