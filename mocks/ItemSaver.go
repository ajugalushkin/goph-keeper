// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	mock "github.com/stretchr/testify/mock"
)

// ItemSaver is an autogenerated mock type for the ItemSaver type
type ItemSaver struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *ItemSaver) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	ret := _m.Called(ctx, item)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Item) (*models.Item, error)); ok {
		return rf(ctx, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Item) *models.Item); ok {
		r0 = rf(ctx, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Item) error); ok {
		r1 = rf(ctx, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, item
func (_m *ItemSaver) Delete(ctx context.Context, item *models.Item) error {
	ret := _m.Called(ctx, item)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Item) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, item
func (_m *ItemSaver) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	ret := _m.Called(ctx, item)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *models.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Item) (*models.Item, error)); ok {
		return rf(ctx, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Item) *models.Item); ok {
		r0 = rf(ctx, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Item) error); ok {
		r1 = rf(ctx, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewItemSaver creates a new instance of ItemSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewItemSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *ItemSaver {
	mock := &ItemSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
