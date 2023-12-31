// Code generated by MockGen. DO NOT EDIT.
// Source: internal/reader/repo/repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/reader/repo/repo.go -package=reader Repo
//

// Package reader is a generated GoMock package.
package reader

import (
	context "context"
	reflect "reflect"

	entity "github.com/bow/neon/internal/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Backend mocks base method.
func (m *MockRepo) Backend() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Backend")
	ret0, _ := ret[0].(string)
	return ret0
}

// Backend indicates an expected call of Backend.
func (mr *MockRepoMockRecorder) Backend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Backend", reflect.TypeOf((*MockRepo)(nil).Backend))
}

// GetStats mocks base method.
func (m *MockRepo) GetStats(arg0 context.Context) (<-chan *entity.Stats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStats", arg0)
	ret0, _ := ret[0].(<-chan *entity.Stats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStats indicates an expected call of GetStats.
func (mr *MockRepoMockRecorder) GetStats(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStats", reflect.TypeOf((*MockRepo)(nil).GetStats), arg0)
}

// ListFeeds mocks base method.
func (m *MockRepo) ListFeeds(arg0 context.Context) (<-chan *entity.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeeds", arg0)
	ret0, _ := ret[0].(<-chan *entity.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeeds indicates an expected call of ListFeeds.
func (mr *MockRepoMockRecorder) ListFeeds(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeeds", reflect.TypeOf((*MockRepo)(nil).ListFeeds), arg0)
}

// PullFeeds mocks base method.
func (m *MockRepo) PullFeeds(arg0 context.Context) (<-chan *entity.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullFeeds", arg0)
	ret0, _ := ret[0].(<-chan *entity.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PullFeeds indicates an expected call of PullFeeds.
func (mr *MockRepoMockRecorder) PullFeeds(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullFeeds", reflect.TypeOf((*MockRepo)(nil).PullFeeds), arg0)
}
