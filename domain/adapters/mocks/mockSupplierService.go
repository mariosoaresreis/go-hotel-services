// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mariosoaresreis/go-hotel/domain/ports/interfaces (interfaces: SupplierService)

// Package services is a generated GoMock package.
package services

import (
	reflect "reflect"
	sync "sync"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/mariosoaresreis/go-hotel/domain"
	interfaces "github.com/mariosoaresreis/go-hotel/domain/ports/interfaces"
	errors "github.com/mariosoaresreis/go-hotel/errors"
)

// MockSupplierService is a mock of SupplierService interface.
type MockSupplierService struct {
	ctrl     *gomock.Controller
	recorder *MockSupplierServiceMockRecorder
}

// MockSupplierServiceMockRecorder is the mock recorder for MockSupplierService.
type MockSupplierServiceMockRecorder struct {
	mock *MockSupplierService
}

// NewMockSupplierService creates a new mock instance.
func NewMockSupplierService(ctrl *gomock.Controller) *MockSupplierService {
	mock := &MockSupplierService{ctrl: ctrl}
	mock.recorder = &MockSupplierServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSupplierService) EXPECT() *MockSupplierServiceMockRecorder {
	return m.recorder
}

// CreateJob mocks base method.
func (m *MockSupplierService) CreateJob(arg0 *domain.Job, arg1 *sync.WaitGroup) (*domain.JobResponse, *errors.ApplicationError) {
	if (arg1 != nil){
		defer arg1.Done()
	}

	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJob", arg0, arg1)
	ret0, _ := ret[0].(*domain.JobResponse)
	ret1, _ := ret[1].(*errors.ApplicationError)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob.
func (mr *MockSupplierServiceMockRecorder) CreateJob(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockSupplierService)(nil).CreateJob), arg0, arg1)
}

// GetAllLocations mocks base method.
func (m *MockSupplierService) GetAllLocations() ([]domain.Location, *errors.ApplicationError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllLocations")
	ret0, _ := ret[0].([]domain.Location)
	ret1, _ := ret[1].(*errors.ApplicationError)
	return ret0, ret1
}

// GetAllLocations indicates an expected call of GetAllLocations.
func (mr *MockSupplierServiceMockRecorder) GetAllLocations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllLocations", reflect.TypeOf((*MockSupplierService)(nil).GetAllLocations))
}

// GetDepartment mocks base method.
func (m *MockSupplierService) GetDepartment(arg0 int64, arg1 *domain.Department, arg2 *sync.WaitGroup) (*domain.Department, *errors.ApplicationError) {
	if (arg2 != nil){
		defer arg2.Done()
	}

	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDepartment", arg0, arg1, arg2)
	ret0, _ := ret[0].(*domain.Department)
	ret1, _ := ret[1].(*errors.ApplicationError)
	arg1.ID = 1
	return ret0, ret1
}

// GetDepartment indicates an expected call of GetDepartment.
func (mr *MockSupplierServiceMockRecorder) GetDepartment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDepartment", reflect.TypeOf((*MockSupplierService)(nil).GetDepartment), arg0, arg1, arg2)
}

// GetJobItem mocks base method.
func (m *MockSupplierService) GetJobItem(arg0 int64, arg1 *domain.JobItemResponse, arg2 *sync.WaitGroup) (*domain.JobItemResponse, *errors.ApplicationError) {
	if (arg2 != nil){
		defer arg2.Done()
	}

	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobItem", arg0, arg1, arg2)
	ret0, _ := ret[0].(*domain.JobItemResponse)
	ret1, _ := ret[1].(*errors.ApplicationError)
	arg1.Id = 1
	return ret0, ret1
}

// GetJobItem indicates an expected call of GetJobItem.
func (mr *MockSupplierServiceMockRecorder) GetJobItem(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobItem", reflect.TypeOf((*MockSupplierService)(nil).GetJobItem), arg0, arg1, arg2)
}

// GetLocation mocks base method.
func (m *MockSupplierService) GetLocation(arg0 int64, arg1 *domain.Location, arg2 *sync.WaitGroup, arg3 interfaces.MessageMap,
	arg4 interfaces.LocationMap) (*domain.Location, *errors.ApplicationError) {
	if (arg2 != nil){
		defer arg2.Done()
	}
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocation", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*domain.Location)
	ret1, _ := ret[1].(*errors.ApplicationError)
	arg1.Id = 1
	arg1.Id = ret0.Id
	arg1.Name = ret0.Name
	arg1.DisplayName = ret0.DisplayName

	arg3.M = make(map[int64]string, 1)
	arg4.M[1] = domain.Location{Id: 2, DisplayName: "displayName", Name: "Name", LocationType: domain.LocationType{ Id: 2, DisplayName: "Floor"}}

	return ret0, ret1
}

// GetLocation indicates an expected call of GetLocation.
func (mr *MockSupplierServiceMockRecorder) GetLocation(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocation", reflect.TypeOf((*MockSupplierService)(nil).GetLocation), arg0, arg1, arg2, arg3, arg4)
}