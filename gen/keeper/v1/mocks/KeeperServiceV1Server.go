// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// KeeperServiceV1Server is an autogenerated mock type for the KeeperServiceV1Server type
type KeeperServiceV1Server struct {
	mock.Mock
}

// CreateItemStreamV1 provides a mock function with given fields: _a0
func (_m *KeeperServiceV1Server) CreateItemStreamV1(_a0 grpc.ClientStreamingServer[v1.CreateItemStreamRequestV1, v1.CreateItemResponseV1]) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateItemStreamV1")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(grpc.ClientStreamingServer[v1.CreateItemStreamRequestV1, v1.CreateItemResponseV1]) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateItemV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) CreateItemV1(_a0 context.Context, _a1 *v1.CreateItemRequestV1) (*v1.CreateItemResponseV1, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateItemV1")
	}

	var r0 *v1.CreateItemResponseV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CreateItemRequestV1) (*v1.CreateItemResponseV1, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CreateItemRequestV1) *v1.CreateItemResponseV1); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.CreateItemResponseV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.CreateItemRequestV1) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteItemV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) DeleteItemV1(_a0 context.Context, _a1 *v1.DeleteItemRequestV1) (*v1.DeleteItemResponseV1, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteItemV1")
	}

	var r0 *v1.DeleteItemResponseV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeleteItemRequestV1) (*v1.DeleteItemResponseV1, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeleteItemRequestV1) *v1.DeleteItemResponseV1); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.DeleteItemResponseV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.DeleteItemRequestV1) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItemStreamV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) GetItemStreamV1(_a0 *v1.GetItemRequestV1, _a1 grpc.ServerStreamingServer[v1.GetItemStreamResponseV1]) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetItemStreamV1")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.GetItemRequestV1, grpc.ServerStreamingServer[v1.GetItemStreamResponseV1]) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetItemV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) GetItemV1(_a0 context.Context, _a1 *v1.GetItemRequestV1) (*v1.GetItemResponseV1, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetItemV1")
	}

	var r0 *v1.GetItemResponseV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.GetItemRequestV1) (*v1.GetItemResponseV1, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.GetItemRequestV1) *v1.GetItemResponseV1); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.GetItemResponseV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.GetItemRequestV1) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListItemsV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) ListItemsV1(_a0 context.Context, _a1 *v1.ListItemsRequestV1) (*v1.ListItemsResponseV1, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListItemsV1")
	}

	var r0 *v1.ListItemsResponseV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ListItemsRequestV1) (*v1.ListItemsResponseV1, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ListItemsRequestV1) *v1.ListItemsResponseV1); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ListItemsResponseV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ListItemsRequestV1) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateItemV1 provides a mock function with given fields: _a0, _a1
func (_m *KeeperServiceV1Server) UpdateItemV1(_a0 context.Context, _a1 *v1.UpdateItemRequestV1) (*v1.UpdateItemResponseV1, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateItemV1")
	}

	var r0 *v1.UpdateItemResponseV1
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.UpdateItemRequestV1) (*v1.UpdateItemResponseV1, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.UpdateItemRequestV1) *v1.UpdateItemResponseV1); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.UpdateItemResponseV1)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.UpdateItemRequestV1) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedKeeperServiceV1Server provides a mock function with given fields:
func (_m *KeeperServiceV1Server) mustEmbedUnimplementedKeeperServiceV1Server() {
	_m.Called()
}

// NewKeeperServiceV1Server creates a new instance of KeeperServiceV1Server. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKeeperServiceV1Server(t interface {
	mock.TestingT
	Cleanup(func())
}) *KeeperServiceV1Server {
	mock := &KeeperServiceV1Server{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
