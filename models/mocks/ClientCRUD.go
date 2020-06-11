// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	models "authserver/models"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ClientCRUD is an autogenerated mock type for the ClientCRUD type
type ClientCRUD struct {
	mock.Mock
}

// GetClientByID provides a mock function with given fields: ID
func (_m *ClientCRUD) GetClientByID(ID uuid.UUID) (*models.Client, error) {
	ret := _m.Called(ID)

	var r0 *models.Client
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Client); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
