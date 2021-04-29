// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cage1016/ms-sample/internal/app/add/service (interfaces: AddService)

// Package automocks is a generated GoMock package.
package automocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAddService is a mock of AddService interface
type MockAddService struct {
	ctrl     *gomock.Controller
	recorder *MockAddServiceMockRecorder
}

// MockAddServiceMockRecorder is the mock recorder for MockAddService
type MockAddServiceMockRecorder struct {
	mock *MockAddService
}

// NewMockAddService creates a new mock instance
func NewMockAddService(ctrl *gomock.Controller) *MockAddService {
	mock := &MockAddService{ctrl: ctrl}
	mock.recorder = &MockAddServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAddService) EXPECT() *MockAddServiceMockRecorder {
	return m.recorder
}

// Sum mocks base method
func (m *MockAddService) Sum(arg0 context.Context, arg1, arg2 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sum", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sum indicates an expected call of Sum
func (mr *MockAddServiceMockRecorder) Sum(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sum", reflect.TypeOf((*MockAddService)(nil).Sum), arg0, arg1, arg2)
}
