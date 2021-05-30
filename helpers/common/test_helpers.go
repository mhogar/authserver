package commonhelpers

import (
	"github.com/stretchr/testify/suite"
)

// AssertContainsSubstrings assets the provided str contains all the expected substrings
func AssertContainsSubstrings(suite *suite.Suite, str string, expectedSubStrs ...string) {
	for _, expectedSubStr := range expectedSubStrs {
		suite.Contains(str, expectedSubStr)
	}
}

// AssertError asserts that err is not nil and its message contains the expects sub strings.
func AssertError(suite *suite.Suite, err error, expectedSubStrs ...string) {
	suite.Require().Error(err)
	AssertContainsSubstrings(suite, err.Error(), expectedSubStrs...)
}

// AssertInternalError assets that err is not nil and its message is an internal error message
func AssertInternalError(suite *suite.Suite, err error) {
	AssertError(suite, err, "internal error")
}
