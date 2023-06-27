// Code generated by mockery v2.28.2. DO NOT EDIT.

package ca_chain

import mock "github.com/stretchr/testify/mock"

// mockFetcher is an autogenerated mock type for the fetcher type
type mockFetcher struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: url
func (_m *mockFetcher) Fetch(url string) ([]byte, error) {
	ret := _m.Called(url)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockFetcher interface {
	mock.TestingT
	Cleanup(func())
}

// newMockFetcher creates a new instance of mockFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockFetcher(t mockConstructorTestingTnewMockFetcher) *mockFetcher {
	mock := &mockFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
