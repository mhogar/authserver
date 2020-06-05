package router_test

import (
	controllermocks "authserver/controllers/mocks"
	"authserver/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
	RequestHandler controllermocks.RequestHandler
	Router         *httprouter.Router
}

func (suite *RouterTestSuite) SetupTest() {
	suite.RequestHandler = controllermocks.RequestHandler{}
	suite.Router = router.CreateRouter(&suite.RequestHandler)
}

func (suite *RouterTestSuite) TestRouter_SendsInternalServerErrorOnPanic() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/user", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("PostAPIUser", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.EqualValues(http.StatusInternalServerError, res.StatusCode)
}

func (suite *RouterTestSuite) TestRouter_PostAPIUserHandledByCorrectHandleFunction() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/user", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("PostAPIUser", mock.Anything, mock.Anything, mock.Anything)

	//act
	_, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.RequestHandler.AssertCalled(suite.T(), "PostAPIUser", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *RouterTestSuite) TestRouter_DeleteAPIUserHandledByCorrectHandleFunction() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/user/1", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("DeleteAPIUser", mock.Anything, mock.Anything, mock.Anything)

	//act
	_, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.RequestHandler.AssertCalled(suite.T(), "DeleteAPIUser", mock.Anything, mock.Anything, mock.MatchedBy(func(params httprouter.Params) bool {
		return params.ByName("id") == "1"
	}))
}

func (suite *RouterTestSuite) TestRouter_PatchAPIUserPasswordHandledByCorrectHandleFunction() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPatch, server.URL+"/api/user/password", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("PatchAPIUserPassword", mock.Anything, mock.Anything, mock.Anything)

	//act
	_, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.RequestHandler.AssertCalled(suite.T(), "PatchAPIUserPassword", mock.Anything, mock.Anything, mock.Anything)
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{})
}
