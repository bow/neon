// Code generated by MockGen. DO NOT EDIT.
// Source: internal/reader/ui/operator.go
//
// Generated by this command:
//
//	mockgen -source=internal/reader/ui/operator.go -package=reader Operator
//

// Package reader is a generated GoMock package.
package reader

import (
	reflect "reflect"

	entity "github.com/bow/neon/internal/entity"
	ui "github.com/bow/neon/internal/reader/ui"
	gomock "go.uber.org/mock/gomock"
)

// MockOperator is a mock of Operator interface.
type MockOperator struct {
	ctrl     *gomock.Controller
	recorder *MockOperatorMockRecorder
}

// MockOperatorMockRecorder is the mock recorder for MockOperator.
type MockOperatorMockRecorder struct {
	mock *MockOperator
}

// NewMockOperator creates a new mock instance.
func NewMockOperator(ctrl *gomock.Controller) *MockOperator {
	mock := &MockOperator{ctrl: ctrl}
	mock.recorder = &MockOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperator) EXPECT() *MockOperatorMockRecorder {
	return m.recorder
}

// ClearStatusBar mocks base method.
func (m *MockOperator) ClearStatusBar(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearStatusBar", arg0)
}

// ClearStatusBar indicates an expected call of ClearStatusBar.
func (mr *MockOperatorMockRecorder) ClearStatusBar(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearStatusBar", reflect.TypeOf((*MockOperator)(nil).ClearStatusBar), arg0)
}

// FocusEntriesPane mocks base method.
func (m *MockOperator) FocusEntriesPane(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FocusEntriesPane", arg0)
}

// FocusEntriesPane indicates an expected call of FocusEntriesPane.
func (mr *MockOperatorMockRecorder) FocusEntriesPane(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FocusEntriesPane", reflect.TypeOf((*MockOperator)(nil).FocusEntriesPane), arg0)
}

// FocusFeedsPane mocks base method.
func (m *MockOperator) FocusFeedsPane(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FocusFeedsPane", arg0)
}

// FocusFeedsPane indicates an expected call of FocusFeedsPane.
func (mr *MockOperatorMockRecorder) FocusFeedsPane(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FocusFeedsPane", reflect.TypeOf((*MockOperator)(nil).FocusFeedsPane), arg0)
}

// FocusNextPane mocks base method.
func (m *MockOperator) FocusNextPane(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FocusNextPane", arg0)
}

// FocusNextPane indicates an expected call of FocusNextPane.
func (mr *MockOperatorMockRecorder) FocusNextPane(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FocusNextPane", reflect.TypeOf((*MockOperator)(nil).FocusNextPane), arg0)
}

// FocusPreviousPane mocks base method.
func (m *MockOperator) FocusPreviousPane(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FocusPreviousPane", arg0)
}

// FocusPreviousPane indicates an expected call of FocusPreviousPane.
func (mr *MockOperatorMockRecorder) FocusPreviousPane(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FocusPreviousPane", reflect.TypeOf((*MockOperator)(nil).FocusPreviousPane), arg0)
}

// FocusReadingPane mocks base method.
func (m *MockOperator) FocusReadingPane(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FocusReadingPane", arg0)
}

// FocusReadingPane indicates an expected call of FocusReadingPane.
func (mr *MockOperatorMockRecorder) FocusReadingPane(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FocusReadingPane", reflect.TypeOf((*MockOperator)(nil).FocusReadingPane), arg0)
}

// GetCurrentFeed mocks base method.
func (m *MockOperator) GetCurrentFeed(arg0 *ui.Display) *entity.Feed {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentFeed", arg0)
	ret0, _ := ret[0].(*entity.Feed)
	return ret0
}

// GetCurrentFeed indicates an expected call of GetCurrentFeed.
func (mr *MockOperatorMockRecorder) GetCurrentFeed(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentFeed", reflect.TypeOf((*MockOperator)(nil).GetCurrentFeed), arg0)
}

// PopulateFeedsPane mocks base method.
func (m *MockOperator) PopulateFeedsPane(arg0 *ui.Display, arg1 func() ([]*entity.Feed, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PopulateFeedsPane", arg0, arg1)
}

// PopulateFeedsPane indicates an expected call of PopulateFeedsPane.
func (mr *MockOperatorMockRecorder) PopulateFeedsPane(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopulateFeedsPane", reflect.TypeOf((*MockOperator)(nil).PopulateFeedsPane), arg0, arg1)
}

// RefreshFeeds mocks base method.
func (m *MockOperator) RefreshFeeds(arg0 *ui.Display, arg1 func() (<-chan entity.PullResult, error), arg2 *entity.Feed) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RefreshFeeds", arg0, arg1, arg2)
}

// RefreshFeeds indicates an expected call of RefreshFeeds.
func (mr *MockOperatorMockRecorder) RefreshFeeds(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshFeeds", reflect.TypeOf((*MockOperator)(nil).RefreshFeeds), arg0, arg1, arg2)
}

// RefreshStats mocks base method.
func (m *MockOperator) RefreshStats(arg0 *ui.Display, arg1 func() (*entity.Stats, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RefreshStats", arg0, arg1)
}

// RefreshStats indicates an expected call of RefreshStats.
func (mr *MockOperatorMockRecorder) RefreshStats(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshStats", reflect.TypeOf((*MockOperator)(nil).RefreshStats), arg0, arg1)
}

// ShowIntroPopup mocks base method.
func (m *MockOperator) ShowIntroPopup(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowIntroPopup", arg0)
}

// ShowIntroPopup indicates an expected call of ShowIntroPopup.
func (mr *MockOperatorMockRecorder) ShowIntroPopup(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowIntroPopup", reflect.TypeOf((*MockOperator)(nil).ShowIntroPopup), arg0)
}

// ToggleAboutPopup mocks base method.
func (m *MockOperator) ToggleAboutPopup(arg0 *ui.Display, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleAboutPopup", arg0, arg1)
}

// ToggleAboutPopup indicates an expected call of ToggleAboutPopup.
func (mr *MockOperatorMockRecorder) ToggleAboutPopup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleAboutPopup", reflect.TypeOf((*MockOperator)(nil).ToggleAboutPopup), arg0, arg1)
}

// ToggleAllFeedsFold mocks base method.
func (m *MockOperator) ToggleAllFeedsFold(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleAllFeedsFold", arg0)
}

// ToggleAllFeedsFold indicates an expected call of ToggleAllFeedsFold.
func (mr *MockOperatorMockRecorder) ToggleAllFeedsFold(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleAllFeedsFold", reflect.TypeOf((*MockOperator)(nil).ToggleAllFeedsFold), arg0)
}

// ToggleCurrentFeedFold mocks base method.
func (m *MockOperator) ToggleCurrentFeedFold(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleCurrentFeedFold", arg0)
}

// ToggleCurrentFeedFold indicates an expected call of ToggleCurrentFeedFold.
func (mr *MockOperatorMockRecorder) ToggleCurrentFeedFold(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleCurrentFeedFold", reflect.TypeOf((*MockOperator)(nil).ToggleCurrentFeedFold), arg0)
}

// ToggleHelpPopup mocks base method.
func (m *MockOperator) ToggleHelpPopup(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleHelpPopup", arg0)
}

// ToggleHelpPopup indicates an expected call of ToggleHelpPopup.
func (mr *MockOperatorMockRecorder) ToggleHelpPopup(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleHelpPopup", reflect.TypeOf((*MockOperator)(nil).ToggleHelpPopup), arg0)
}

// ToggleStatsPopup mocks base method.
func (m *MockOperator) ToggleStatsPopup(arg0 *ui.Display, arg1 func() (*entity.Stats, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleStatsPopup", arg0, arg1)
}

// ToggleStatsPopup indicates an expected call of ToggleStatsPopup.
func (mr *MockOperatorMockRecorder) ToggleStatsPopup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleStatsPopup", reflect.TypeOf((*MockOperator)(nil).ToggleStatsPopup), arg0, arg1)
}

// ToggleStatusBar mocks base method.
func (m *MockOperator) ToggleStatusBar(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleStatusBar", arg0)
}

// ToggleStatusBar indicates an expected call of ToggleStatusBar.
func (mr *MockOperatorMockRecorder) ToggleStatusBar(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleStatusBar", reflect.TypeOf((*MockOperator)(nil).ToggleStatusBar), arg0)
}

// UnfocusFront mocks base method.
func (m *MockOperator) UnfocusFront(arg0 *ui.Display) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnfocusFront", arg0)
}

// UnfocusFront indicates an expected call of UnfocusFront.
func (mr *MockOperatorMockRecorder) UnfocusFront(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnfocusFront", reflect.TypeOf((*MockOperator)(nil).UnfocusFront), arg0)
}
