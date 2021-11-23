// Code generated by MockGen. DO NOT EDIT.
// Source: spenmo/payment-processing/payment-processing/internal/pkg/store/mysql (interfaces: LimitStore)

// Package storemocks is a generated GoMock package.
package storemocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
	reflect "reflect"
	store "spenmo/payment-processing/payment-processing/internal/pkg/store"
	mysql "spenmo/payment-processing/payment-processing/internal/pkg/store/mysql"
)

// MockLimitStore is a mock of LimitStore interface
type MockLimitStore struct {
	ctrl     *gomock.Controller
	recorder *MockLimitStoreMockRecorder
}

// MockLimitStoreMockRecorder is the mock recorder for MockLimitStore
type MockLimitStoreMockRecorder struct {
	mock *MockLimitStore
}

// NewMockLimitStore creates a new mock instance
func NewMockLimitStore(ctrl *gomock.Controller) *MockLimitStore {
	mock := &MockLimitStore{ctrl: ctrl}
	mock.recorder = &MockLimitStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLimitStore) EXPECT() *MockLimitStoreMockRecorder {
	return m.recorder
}

// BeginX mocks base method
func (m *MockLimitStore) BeginX() (*sqlx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginX")
	ret0, _ := ret[0].(*sqlx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginX indicates an expected call of BeginX
func (mr *MockLimitStoreMockRecorder) BeginX() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginX", reflect.TypeOf((*MockLimitStore)(nil).BeginX))
}

// CommitX mocks base method
func (m *MockLimitStore) CommitX(arg0 *sqlx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitX", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitX indicates an expected call of CommitX
func (mr *MockLimitStoreMockRecorder) CommitX(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitX", reflect.TypeOf((*MockLimitStore)(nil).CommitX), arg0)
}

// CreateLimit mocks base method
func (m *MockLimitStore) CreateLimit(arg0 context.Context, arg1 mysql.Execer, arg2 *store.Limit) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLimit", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLimit indicates an expected call of CreateLimit
func (mr *MockLimitStoreMockRecorder) CreateLimit(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLimit", reflect.TypeOf((*MockLimitStore)(nil).CreateLimit), arg0, arg1, arg2)
}

// DBX mocks base method
func (m *MockLimitStore) DBX() *sqlx.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBX")
	ret0, _ := ret[0].(*sqlx.DB)
	return ret0
}

// DBX indicates an expected call of DBX
func (mr *MockLimitStoreMockRecorder) DBX() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBX", reflect.TypeOf((*MockLimitStore)(nil).DBX))
}

// GetLimits mocks base method
func (m *MockLimitStore) GetLimits(arg0 context.Context, arg1 mysql.Querier, arg2 int, arg3 int64) ([]*store.Limit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimits", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*store.Limit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLimits indicates an expected call of GetLimits
func (mr *MockLimitStoreMockRecorder) GetLimits(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimits", reflect.TypeOf((*MockLimitStore)(nil).GetLimits), arg0, arg1, arg2, arg3)
}

// RollbackX mocks base method
func (m *MockLimitStore) RollbackX(arg0 *sqlx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackX", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollbackX indicates an expected call of RollbackX
func (mr *MockLimitStoreMockRecorder) RollbackX(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackX", reflect.TypeOf((*MockLimitStore)(nil).RollbackX), arg0)
}