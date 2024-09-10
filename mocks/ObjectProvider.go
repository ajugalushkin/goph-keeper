// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/ajugalushkin/goph-keeper/server/internal/dto/models"
	mock "github.com/stretchr/testify/mock"
)

// ObjectProvider is an autogenerated mock type for the ObjectProvider type
type ObjectProvider struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, objectID
func (_m *ObjectProvider) Get(ctx context.Context, objectID string) (*models.File, error) {
	ret := _m.Called(ctx, objectID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.File
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.File, error)); ok {
		return rf(ctx, objectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.File); ok {
		r0 = rf(ctx, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.File)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewObjectProvider creates a new instance of ObjectProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObjectProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObjectProvider {
	mock := &ObjectProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
