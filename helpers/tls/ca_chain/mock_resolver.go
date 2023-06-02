// Code generated by mockery v2.28.2. DO NOT EDIT.

package ca_chain

import (
	x509 "crypto/x509"

	mock "github.com/stretchr/testify/mock"
)

// mockResolver is an autogenerated mock type for the resolver type
type mockResolver struct {
	mock.Mock
}

// Resolve provides a mock function with given fields: certs
func (_m *mockResolver) Resolve(certs []*x509.Certificate) ([]*x509.Certificate, error) {
	ret := _m.Called(certs)

	var r0 []*x509.Certificate
	var r1 error
	if rf, ok := ret.Get(0).(func([]*x509.Certificate) ([]*x509.Certificate, error)); ok {
		return rf(certs)
	}
	if rf, ok := ret.Get(0).(func([]*x509.Certificate) []*x509.Certificate); ok {
		r0 = rf(certs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*x509.Certificate)
		}
	}

	if rf, ok := ret.Get(1).(func([]*x509.Certificate) error); ok {
		r1 = rf(certs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockResolver interface {
	mock.TestingT
	Cleanup(func())
}

// newMockResolver creates a new instance of mockResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockResolver(t mockConstructorTestingTnewMockResolver) *mockResolver {
	mock := &mockResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
