// Code generated by mockery v2.28.2. DO NOT EDIT.

package vault

import mock "github.com/stretchr/testify/mock"

// MockAuthMethod is an autogenerated mock type for the AuthMethod type
type MockAuthMethod struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: client
func (_m *MockAuthMethod) Authenticate(client Client) error {
	ret := _m.Called(client)

	var r0 error
	if rf, ok := ret.Get(0).(func(Client) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *MockAuthMethod) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Token provides a mock function with given fields:
func (_m *MockAuthMethod) Token() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewMockAuthMethod interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAuthMethod creates a new instance of MockAuthMethod. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAuthMethod(t mockConstructorTestingTNewMockAuthMethod) *MockAuthMethod {
	mock := &MockAuthMethod{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
