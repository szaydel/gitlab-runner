// Code generated by mockery v2.28.2. DO NOT EDIT.

package user

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockInspect is an autogenerated mock type for the Inspect type
type MockInspect struct {
	mock.Mock
}

// GID provides a mock function with given fields: ctx, containerID
func (_m *MockInspect) GID(ctx context.Context, containerID string) (int, error) {
	ret := _m.Called(ctx, containerID)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, containerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, containerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, containerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsRoot provides a mock function with given fields: ctx, imageID
func (_m *MockInspect) IsRoot(ctx context.Context, imageID string) (bool, error) {
	ret := _m.Called(ctx, imageID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, imageID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, imageID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, imageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UID provides a mock function with given fields: ctx, containerID
func (_m *MockInspect) UID(ctx context.Context, containerID string) (int, error) {
	ret := _m.Called(ctx, containerID)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, containerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, containerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, containerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockInspect interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockInspect creates a new instance of MockInspect. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockInspect(t mockConstructorTestingTNewMockInspect) *MockInspect {
	mock := &MockInspect{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
