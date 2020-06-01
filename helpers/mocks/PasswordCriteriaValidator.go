// Code generated by mockery v1.1.2. DO NOT EDIT.

package mocks

import (
	helpers "authserver/helpers"

	mock "github.com/stretchr/testify/mock"
)

// PasswordCriteriaValidator is an autogenerated mock type for the PasswordCriteriaValidator type
type PasswordCriteriaValidator struct {
	mock.Mock
}

// ValidatePasswordCriteria provides a mock function with given fields: password
func (_m *PasswordCriteriaValidator) ValidatePasswordCriteria(password string) helpers.ValidatePasswordCriteriaError {
	ret := _m.Called(password)

	var r0 helpers.ValidatePasswordCriteriaError
	if rf, ok := ret.Get(0).(func(string) helpers.ValidatePasswordCriteriaError); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(helpers.ValidatePasswordCriteriaError)
	}

	return r0
}
