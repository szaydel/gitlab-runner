// Code generated by mockery v2.28.2. DO NOT EDIT.

package vault

import mock "github.com/stretchr/testify/mock"

// mockApiClient is an autogenerated mock type for the apiClient type
type mockApiClient struct {
	mock.Mock
}

// Logical provides a mock function with given fields:
func (_m *mockApiClient) Logical() apiClientLogical {
	ret := _m.Called()

	var r0 apiClientLogical
	if rf, ok := ret.Get(0).(func() apiClientLogical); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(apiClientLogical)
		}
	}

	return r0
}

// SetNamespace provides a mock function with given fields: ns
func (_m *mockApiClient) SetNamespace(ns string) {
	_m.Called(ns)
}

// SetToken provides a mock function with given fields: v
func (_m *mockApiClient) SetToken(v string) {
	_m.Called(v)
}

// Sys provides a mock function with given fields:
func (_m *mockApiClient) Sys() apiClientSys {
	ret := _m.Called()

	var r0 apiClientSys
	if rf, ok := ret.Get(0).(func() apiClientSys); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(apiClientSys)
		}
	}

	return r0
}

type mockConstructorTestingTnewMockApiClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockApiClient creates a new instance of mockApiClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockApiClient(t mockConstructorTestingTnewMockApiClient) *mockApiClient {
	mock := &mockApiClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
