package router_test

import (
	controllermocks "authserver/controllers/mocks"
	databasemocks "authserver/database/mocks"
	"authserver/router"
	"authserver/router/mocks"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
	ControllersMock   controllermocks.Controllers
	AuthenticatorMock mocks.Authenticator
	Transaction       databasemocks.Transaction
	Router            *httprouter.Router
}

func (suite *RouterTestSuite) SetupTest() {
	suite.ControllersMock = controllermocks.Controllers{}
	suite.AuthenticatorMock = mocks.Authenticator{}
	suite.Transaction = databasemocks.Transaction{}

	tf := databasemocks.TransactionFactory{}
	tf.On("CreateTransaction").Return(&suite.Transaction, nil)

	suite.Transaction.On("CommitTransaction").Return(nil)
	suite.Transaction.On("RollbackTransaction")

	rf := router.RouterFactory{
		Controllers:        &suite.ControllersMock,
		Authenticator:      &suite.AuthenticatorMock,
		TransactionFactory: &tf,
	}
	suite.Router = rf.CreateRouter()
}
