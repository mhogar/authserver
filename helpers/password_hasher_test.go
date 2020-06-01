package helpers_test

import (
	"authserver/helpers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BCryptPasswordHasherTestSuite struct {
	suite.Suite
	BCryptPasswordHasher helpers.BCryptPasswordHasher
}

func (suite *BCryptPasswordHasherTestSuite) SetupTest() {
	suite.BCryptPasswordHasher = helpers.BCryptPasswordHasher{}
}

func (suite *BCryptPasswordHasherTestSuite) TestHashPassword_WithNoError_ReturnsHashAndNilError() {
	hash, err := suite.BCryptPasswordHasher.HashPassword("password")
	suite.NotNil(hash)
	suite.NoError(err)
}

func (suite *BCryptPasswordHasherTestSuite) TestComparePasswords_WherePasswordMatchesHash_ReturnsNilError() {
	//arrange
	password := "password"
	hash, err := suite.BCryptPasswordHasher.HashPassword(password)
	suite.NoError(err)

	//act
	err = suite.BCryptPasswordHasher.ComparePasswords(hash, password)

	//assert
	suite.NoError(err)
}

func (suite *BCryptPasswordHasherTestSuite) TestComparePasswords_WherePasswordDoesNotMatchHash_ReturnsError() {
	//arrange
	password := "password"

	//act
	err := suite.BCryptPasswordHasher.ComparePasswords([]byte("incorrect hash"), password)

	//assert
	suite.Error(err)
}

func TestBCryptPasswordHasherTestSuite(t *testing.T) {
	suite.Run(t, &BCryptPasswordHasherTestSuite{})
}
