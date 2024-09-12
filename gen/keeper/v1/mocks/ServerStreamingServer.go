// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"

	metadata "google.golang.org/grpc/metadata"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// ServerStreamingServer is an autogenerated mock type for the ServerStreamingServer type
type ServerStreamingServer[Res interface{}] struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *ServerStreamingServer[Res]) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	r0:= context.TODO()
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// RecvMsg provides a mock function with given fields: m
func (_m *ServerStreamingServer[Res]) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: _a0
func (_m *ServerStreamingServer[Res]) Send(_a0 *v1.GetItemStreamResponseV1) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.GetItemStreamResponseV1) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendHeader provides a mock function with given fields: _a0
func (_m *ServerStreamingServer[Res]) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMsg provides a mock function with given fields: m
func (_m *ServerStreamingServer[Res]) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHeader provides a mock function with given fields: _a0
func (_m *ServerStreamingServer[Res]) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *ServerStreamingServer[Res]) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// NewServerStreamingServer creates a new instance of ServerStreamingServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServerStreamingServer[Res interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *ServerStreamingServer[Res] {
	mock := &ServerStreamingServer[Res]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
