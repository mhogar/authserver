package helpers_test

import (
	"authserver/config"
	"authserver/helpers"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ConfigPasswordCriteriaValidatorTestSuite struct {
	suite.Suite
	ConfigPasswordCriteriaValidator helpers.ConfigPasswordCriteriaValidator
	Criteria                        config.PasswordCriteriaConfig
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) SetupTest() {
	suite.ConfigPasswordCriteriaValidator = helpers.ConfigPasswordCriteriaValidator{}
	suite.Criteria = config.PasswordCriteriaConfig{
		MinLength:        4,
		RequireLowerCase: false,
		RequireUpperCase: false,
		RequireDigit:     false,
		RequireSymbol:    false,
	}
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) TestValidatePasswordCriteria_MinLengthCriteriaTests() {
	var expectedStatus int
	var password string

	testCase := func() {
		viper.Set("password", suite.Criteria)

		verr := suite.ConfigPasswordCriteriaValidator.ValidatePasswordCriteria(password)
		suite.Equal(expectedStatus, verr.Status)
	}

	password = "aaa"
	expectedStatus = helpers.ValidatePasswordCriteriaTooShort
	suite.Run("PaswordOneLessThanMinLength_ReturnsValidatePasswordCriteriaTooShort", testCase)

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("PaswordEqualToMinLength_ReturnsValidatePasswordCriteriaValid", testCase)

	password = "aaaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("PaswordGreaterThanMinLength_ReturnsValidatePasswordCriteriaValid", testCase)
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) TestValidatePasswordCriteria_RequireLowerCaseCriteriaTests() {
	var expectedStatus int
	var password string

	testCase := func() {
		viper.Set("password", suite.Criteria)

		verr := suite.ConfigPasswordCriteriaValidator.ValidatePasswordCriteria(password)
		suite.Equal(expectedStatus, verr.Status)
	}

	password = "AAAA"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("LowerCaseLetterNotRequiredAndNotContainLowerCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("LowerCaseLetterNotRequiredAndContainsLowerCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)

	suite.Criteria.RequireLowerCase = true

	password = "AAAA"
	expectedStatus = helpers.ValidatePasswordCriteriaMissingLowerCaseLetter
	suite.Run("LowerCaseLetterRequiredAndNotContainLowerCaseLetter_ReturnsValidatePasswordCriteriaMissingLowerCaseLetter", testCase)

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("LowerCaseLetterRequiredAndContainsLowerCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) TestValidatePasswordCriteria_RequireUpperCaseCriteriaTests() {
	var expectedStatus int
	var password string

	testCase := func() {
		viper.Set("password", suite.Criteria)

		verr := suite.ConfigPasswordCriteriaValidator.ValidatePasswordCriteria(password)
		suite.Equal(expectedStatus, verr.Status)
	}

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("UpperCaseLetterNotRequiredAndNotContainUpperCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)

	password = "AAAA"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("UpperCaseLetterNotRequiredAndContainsUpperCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)

	suite.Criteria.RequireUpperCase = true

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaMissingUpperCaseLetter
	suite.Run("UpperCaseLetterRequiredAndNotContainUpperCaseLetter_ReturnsValidatePasswordCriteriaMissingUpperCaseLetter", testCase)

	password = "AAAA"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("UpperCaseLetterRequiredAndContainsUpperCaseLetter_ReturnsValidatePasswordCriteriaValid", testCase)
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) TestValidatePasswordCriteria_RequireDigitCriteriaTests() {
	var expectedStatus int
	var password string

	testCase := func() {
		viper.Set("password", suite.Criteria)

		verr := suite.ConfigPasswordCriteriaValidator.ValidatePasswordCriteria(password)
		suite.Equal(expectedStatus, verr.Status)
	}

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("DigitNotRequiredAndDoesNotContainDigit_ReturnsValidatePasswordCriteriaValid", testCase)

	password = "1234"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("DigitNotRequiredAndContainsDigit_ReturnsValidatePasswordCriteriaValid", testCase)

	suite.Criteria.RequireDigit = true

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaMissingDigit
	suite.Run("DigitRequiredAndNotContainDigit_ReturnsValidatePasswordCriteriaMissingDigit", testCase)

	password = "1234"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("DigitRequiredAndContainsDigit_ReturnsValidatePasswordCriteriaValid", testCase)
}

func (suite *ConfigPasswordCriteriaValidatorTestSuite) TestValidatePasswordCriteria_RequireSymbolCriteriaTests() {
	var expectedStatus int
	var password string

	testCase := func() {
		viper.Set("password", suite.Criteria)

		verr := suite.ConfigPasswordCriteriaValidator.ValidatePasswordCriteria(password)
		suite.Equal(expectedStatus, verr.Status)
	}

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("SymbolNotRequiredAndDoesNotContainSymbol_ReturnsValidatePasswordCriteriaValid", testCase)

	password = "&$%*"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("SymbolNotRequiredAndContainsSymbol_ReturnsValidatePasswordCriteriaValid", testCase)

	suite.Criteria.RequireSymbol = true

	password = "aaaa"
	expectedStatus = helpers.ValidatePasswordCriteriaMissingSymbol
	suite.Run("SymbolRequiredAndNotContainSymbol_ReturnsValidatePasswordCriteriaMissingSymbol", testCase)

	password = "&$%*"
	expectedStatus = helpers.ValidatePasswordCriteriaValid
	suite.Run("SymbolRequiredAndContainsSymbol_ReturnsValidatePasswordCriteriaValid", testCase)
}

func TestConfigPasswordCriteriaValidatorTestSuite(t *testing.T) {
	suite.Run(t, &ConfigPasswordCriteriaValidatorTestSuite{})
}
