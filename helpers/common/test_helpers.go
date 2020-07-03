package commonhelpers

import (
	"github.com/stretchr/testify/suite"
)

// AssertError asserts that err is not nil and its message contains the expects sub strings.
func AssertError(suite *suite.Suite, err error, expectedSubStrs ...string) {
	suite.Require().Error(err)
	for _, expectedSubStr := range expectedSubStrs {
		suite.Contains(err.Error(), expectedSubStr)
	}
}
