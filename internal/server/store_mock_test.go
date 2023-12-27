// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store.go

// Package server is a generated GoMock package.
package server

import (
	context "context"
	reflect "reflect"

	internal "github.com/bow/neon/internal"
	gomock "github.com/golang/mock/gomock"
)

// MockFeedStore is a mock of FeedStore interface.
type MockFeedStore struct {
	ctrl     *gomock.Controller
	recorder *MockFeedStoreMockRecorder
}

// MockFeedStoreMockRecorder is the mock recorder for MockFeedStore.
type MockFeedStoreMockRecorder struct {
	mock *MockFeedStore
}

// NewMockFeedStore creates a new mock instance.
func NewMockFeedStore(ctrl *gomock.Controller) *MockFeedStore {
	mock := &MockFeedStore{ctrl: ctrl}
	mock.recorder = &MockFeedStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeedStore) EXPECT() *MockFeedStoreMockRecorder {
	return m.recorder
}

// AddFeed mocks base method.
func (m *MockFeedStore) AddFeed(ctx context.Context, feedURL string, title, desc *string, tags []string, isStarred *bool) (*internal.Feed, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFeed", ctx, feedURL, title, desc, tags, isStarred)
	ret0, _ := ret[0].(*internal.Feed)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddFeed indicates an expected call of AddFeed.
func (mr *MockFeedStoreMockRecorder) AddFeed(ctx, feedURL, title, desc, tags, isStarred interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFeed", reflect.TypeOf((*MockFeedStore)(nil).AddFeed), ctx, feedURL, title, desc, tags, isStarred)
}

// DeleteFeeds mocks base method.
func (m *MockFeedStore) DeleteFeeds(ctx context.Context, ids []internal.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFeeds", ctx, ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFeeds indicates an expected call of DeleteFeeds.
func (mr *MockFeedStoreMockRecorder) DeleteFeeds(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFeeds", reflect.TypeOf((*MockFeedStore)(nil).DeleteFeeds), ctx, ids)
}

// EditEntries mocks base method.
func (m *MockFeedStore) EditEntries(ctx context.Context, ops []*internal.EntryEditOp) ([]*internal.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditEntries", ctx, ops)
	ret0, _ := ret[0].([]*internal.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditEntries indicates an expected call of EditEntries.
func (mr *MockFeedStoreMockRecorder) EditEntries(ctx, ops interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditEntries", reflect.TypeOf((*MockFeedStore)(nil).EditEntries), ctx, ops)
}

// EditFeeds mocks base method.
func (m *MockFeedStore) EditFeeds(ctx context.Context, ops []*internal.FeedEditOp) ([]*internal.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditFeeds", ctx, ops)
	ret0, _ := ret[0].([]*internal.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditFeeds indicates an expected call of EditFeeds.
func (mr *MockFeedStoreMockRecorder) EditFeeds(ctx, ops interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditFeeds", reflect.TypeOf((*MockFeedStore)(nil).EditFeeds), ctx, ops)
}

// ExportSubscription mocks base method.
func (m *MockFeedStore) ExportSubscription(ctx context.Context, title *string) (*internal.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportSubscription", ctx, title)
	ret0, _ := ret[0].(*internal.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportSubscription indicates an expected call of ExportSubscription.
func (mr *MockFeedStoreMockRecorder) ExportSubscription(ctx, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportSubscription", reflect.TypeOf((*MockFeedStore)(nil).ExportSubscription), ctx, title)
}

// GetEntry mocks base method.
func (m *MockFeedStore) GetEntry(ctx context.Context, id internal.ID) (*internal.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", ctx, id)
	ret0, _ := ret[0].(*internal.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockFeedStoreMockRecorder) GetEntry(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockFeedStore)(nil).GetEntry), ctx, id)
}

// GetGlobalStats mocks base method.
func (m *MockFeedStore) GetGlobalStats(ctx context.Context) (*internal.Stats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGlobalStats", ctx)
	ret0, _ := ret[0].(*internal.Stats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGlobalStats indicates an expected call of GetGlobalStats.
func (mr *MockFeedStoreMockRecorder) GetGlobalStats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGlobalStats", reflect.TypeOf((*MockFeedStore)(nil).GetGlobalStats), ctx)
}

// ImportSubscription mocks base method.
func (m *MockFeedStore) ImportSubscription(ctx context.Context, sub *internal.Subscription) (int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportSubscription", ctx, sub)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ImportSubscription indicates an expected call of ImportSubscription.
func (mr *MockFeedStoreMockRecorder) ImportSubscription(ctx, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportSubscription", reflect.TypeOf((*MockFeedStore)(nil).ImportSubscription), ctx, sub)
}

// ListEntries mocks base method.
func (m *MockFeedStore) ListEntries(ctx context.Context, feedIDs []internal.ID, isBookmarked *bool) ([]*internal.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntries", ctx, feedIDs, isBookmarked)
	ret0, _ := ret[0].([]*internal.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntries indicates an expected call of ListEntries.
func (mr *MockFeedStoreMockRecorder) ListEntries(ctx, feedIDs, isBookmarked interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntries", reflect.TypeOf((*MockFeedStore)(nil).ListEntries), ctx, feedIDs, isBookmarked)
}

// ListFeeds mocks base method.
func (m *MockFeedStore) ListFeeds(ctx context.Context) ([]*internal.Feed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeeds", ctx)
	ret0, _ := ret[0].([]*internal.Feed)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeeds indicates an expected call of ListFeeds.
func (mr *MockFeedStoreMockRecorder) ListFeeds(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeeds", reflect.TypeOf((*MockFeedStore)(nil).ListFeeds), ctx)
}

// PullFeeds mocks base method.
func (m *MockFeedStore) PullFeeds(ctx context.Context, ids []internal.ID) <-chan internal.PullResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullFeeds", ctx, ids)
	ret0, _ := ret[0].(<-chan internal.PullResult)
	return ret0
}

// PullFeeds indicates an expected call of PullFeeds.
func (mr *MockFeedStoreMockRecorder) PullFeeds(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullFeeds", reflect.TypeOf((*MockFeedStore)(nil).PullFeeds), ctx, ids)
}
