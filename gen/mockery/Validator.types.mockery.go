// Code generated by mockery v2.50.4. DO NOT EDIT.

package mockery

import (
	"context"

	mock "github.com/stretchr/testify/mock"
	ast "github.com/walteh/go-tmpl-typer/pkg/ast"
)

// MockValidator_types is an autogenerated mock type for the Validator type
type MockValidator_types struct {
	mock.Mock
}

type MockValidator_types_Expecter struct {
	mock *mock.Mock
}

func (_m *MockValidator_types) EXPECT() *MockValidator_types_Expecter {
	return &MockValidator_types_Expecter{mock: &_m.Mock}
}

// ValidateType provides a mock function with given fields: ctx, typePath, registry
func (_m *MockValidator_types) ValidateType(ctx context.Context, typePath string, registry *ast.TypeRegistry) error {
	ret := _m.Called(ctx, typePath, registry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *ast.TypeRegistry) error); ok {
		r0 = rf(ctx, typePath, registry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockValidator_types_ValidateType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateType'
type MockValidator_types_ValidateType_Call struct {
	*mock.Call
}

// ValidateType is a helper method to define mock.On call
//   - ctx context.Context
//   - typePath string
//   - registry *ast.TypeRegistry
func (_e *MockValidator_types_Expecter) ValidateType(ctx interface{}, typePath interface{}, registry interface{}) *MockValidator_types_ValidateType_Call {
	return &MockValidator_types_ValidateType_Call{Call: _e.mock.On("ValidateType", ctx, typePath, registry)}
}

func (_c *MockValidator_types_ValidateType_Call) Run(run func(ctx context.Context, typePath string, registry *ast.TypeRegistry)) *MockValidator_types_ValidateType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*ast.TypeRegistry))
	})
	return _c
}

func (_c *MockValidator_types_ValidateType_Call) Return(_a0 error) *MockValidator_types_ValidateType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockValidator_types_ValidateType_Call) RunAndReturn(run func(context.Context, string, *ast.TypeRegistry) error) *MockValidator_types_ValidateType_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockValidator_types creates a new instance of MockValidator_types. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockValidator_types(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockValidator_types {
	mock := &MockValidator_types{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}