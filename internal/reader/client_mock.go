// Code generated by MockGen. DO NOT EDIT.
// Source: api/lens_grpc.pb.go

// Package reader is a generated GoMock package.
package reader

import (
	context "context"
	reflect "reflect"

	api "github.com/bow/lens/api"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockLensClient is a mock of LensClient interface.
type MockLensClient struct {
	ctrl     *gomock.Controller
	recorder *MockLensClientMockRecorder
}

// MockLensClientMockRecorder is the mock recorder for MockLensClient.
type MockLensClientMockRecorder struct {
	mock *MockLensClient
}

// NewMockLensClient creates a new mock instance.
func NewMockLensClient(ctrl *gomock.Controller) *MockLensClient {
	mock := &MockLensClient{ctrl: ctrl}
	mock.recorder = &MockLensClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLensClient) EXPECT() *MockLensClientMockRecorder {
	return m.recorder
}

// AddFeed mocks base method.
func (m *MockLensClient) AddFeed(ctx context.Context, in *api.AddFeedRequest, opts ...grpc.CallOption) (*api.AddFeedResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddFeed", varargs...)
	ret0, _ := ret[0].(*api.AddFeedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFeed indicates an expected call of AddFeed.
func (mr *MockLensClientMockRecorder) AddFeed(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFeed", reflect.TypeOf((*MockLensClient)(nil).AddFeed), varargs...)
}

// DeleteFeeds mocks base method.
func (m *MockLensClient) DeleteFeeds(ctx context.Context, in *api.DeleteFeedsRequest, opts ...grpc.CallOption) (*api.DeleteFeedsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFeeds", varargs...)
	ret0, _ := ret[0].(*api.DeleteFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFeeds indicates an expected call of DeleteFeeds.
func (mr *MockLensClientMockRecorder) DeleteFeeds(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFeeds", reflect.TypeOf((*MockLensClient)(nil).DeleteFeeds), varargs...)
}

// EditEntries mocks base method.
func (m *MockLensClient) EditEntries(ctx context.Context, in *api.EditEntriesRequest, opts ...grpc.CallOption) (*api.EditEntriesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditEntries", varargs...)
	ret0, _ := ret[0].(*api.EditEntriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditEntries indicates an expected call of EditEntries.
func (mr *MockLensClientMockRecorder) EditEntries(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditEntries", reflect.TypeOf((*MockLensClient)(nil).EditEntries), varargs...)
}

// EditFeeds mocks base method.
func (m *MockLensClient) EditFeeds(ctx context.Context, in *api.EditFeedsRequest, opts ...grpc.CallOption) (*api.EditFeedsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditFeeds", varargs...)
	ret0, _ := ret[0].(*api.EditFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditFeeds indicates an expected call of EditFeeds.
func (mr *MockLensClientMockRecorder) EditFeeds(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditFeeds", reflect.TypeOf((*MockLensClient)(nil).EditFeeds), varargs...)
}

// ExportOPML mocks base method.
func (m *MockLensClient) ExportOPML(ctx context.Context, in *api.ExportOPMLRequest, opts ...grpc.CallOption) (*api.ExportOPMLResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExportOPML", varargs...)
	ret0, _ := ret[0].(*api.ExportOPMLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportOPML indicates an expected call of ExportOPML.
func (mr *MockLensClientMockRecorder) ExportOPML(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportOPML", reflect.TypeOf((*MockLensClient)(nil).ExportOPML), varargs...)
}

// GetEntry mocks base method.
func (m *MockLensClient) GetEntry(ctx context.Context, in *api.GetEntryRequest, opts ...grpc.CallOption) (*api.GetEntryResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetEntry", varargs...)
	ret0, _ := ret[0].(*api.GetEntryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockLensClientMockRecorder) GetEntry(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockLensClient)(nil).GetEntry), varargs...)
}

// GetInfo mocks base method.
func (m *MockLensClient) GetInfo(ctx context.Context, in *api.GetInfoRequest, opts ...grpc.CallOption) (*api.GetInfoResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetInfo", varargs...)
	ret0, _ := ret[0].(*api.GetInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfo indicates an expected call of GetInfo.
func (mr *MockLensClientMockRecorder) GetInfo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfo", reflect.TypeOf((*MockLensClient)(nil).GetInfo), varargs...)
}

// GetStats mocks base method.
func (m *MockLensClient) GetStats(ctx context.Context, in *api.GetStatsRequest, opts ...grpc.CallOption) (*api.GetStatsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetStats", varargs...)
	ret0, _ := ret[0].(*api.GetStatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStats indicates an expected call of GetStats.
func (mr *MockLensClientMockRecorder) GetStats(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStats", reflect.TypeOf((*MockLensClient)(nil).GetStats), varargs...)
}

// ImportOPML mocks base method.
func (m *MockLensClient) ImportOPML(ctx context.Context, in *api.ImportOPMLRequest, opts ...grpc.CallOption) (*api.ImportOPMLResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ImportOPML", varargs...)
	ret0, _ := ret[0].(*api.ImportOPMLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportOPML indicates an expected call of ImportOPML.
func (mr *MockLensClientMockRecorder) ImportOPML(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportOPML", reflect.TypeOf((*MockLensClient)(nil).ImportOPML), varargs...)
}

// ListEntries mocks base method.
func (m *MockLensClient) ListEntries(ctx context.Context, in *api.ListEntriesRequest, opts ...grpc.CallOption) (*api.ListEntriesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListEntries", varargs...)
	ret0, _ := ret[0].(*api.ListEntriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntries indicates an expected call of ListEntries.
func (mr *MockLensClientMockRecorder) ListEntries(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntries", reflect.TypeOf((*MockLensClient)(nil).ListEntries), varargs...)
}

// ListFeeds mocks base method.
func (m *MockLensClient) ListFeeds(ctx context.Context, in *api.ListFeedsRequest, opts ...grpc.CallOption) (*api.ListFeedsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFeeds", varargs...)
	ret0, _ := ret[0].(*api.ListFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeeds indicates an expected call of ListFeeds.
func (mr *MockLensClientMockRecorder) ListFeeds(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeeds", reflect.TypeOf((*MockLensClient)(nil).ListFeeds), varargs...)
}

// PullFeeds mocks base method.
func (m *MockLensClient) PullFeeds(ctx context.Context, in *api.PullFeedsRequest, opts ...grpc.CallOption) (api.Lens_PullFeedsClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PullFeeds", varargs...)
	ret0, _ := ret[0].(api.Lens_PullFeedsClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PullFeeds indicates an expected call of PullFeeds.
func (mr *MockLensClientMockRecorder) PullFeeds(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullFeeds", reflect.TypeOf((*MockLensClient)(nil).PullFeeds), varargs...)
}

// MockLens_PullFeedsClient is a mock of Lens_PullFeedsClient interface.
type MockLens_PullFeedsClient struct {
	ctrl     *gomock.Controller
	recorder *MockLens_PullFeedsClientMockRecorder
}

// MockLens_PullFeedsClientMockRecorder is the mock recorder for MockLens_PullFeedsClient.
type MockLens_PullFeedsClientMockRecorder struct {
	mock *MockLens_PullFeedsClient
}

// NewMockLens_PullFeedsClient creates a new mock instance.
func NewMockLens_PullFeedsClient(ctrl *gomock.Controller) *MockLens_PullFeedsClient {
	mock := &MockLens_PullFeedsClient{ctrl: ctrl}
	mock.recorder = &MockLens_PullFeedsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLens_PullFeedsClient) EXPECT() *MockLens_PullFeedsClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockLens_PullFeedsClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockLens_PullFeedsClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockLens_PullFeedsClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockLens_PullFeedsClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).Context))
}

// Header mocks base method.
func (m *MockLens_PullFeedsClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockLens_PullFeedsClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockLens_PullFeedsClient) Recv() (*api.PullFeedsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*api.PullFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockLens_PullFeedsClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockLens_PullFeedsClient) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockLens_PullFeedsClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).RecvMsg), m)
}

// SendMsg mocks base method.
func (m_2 *MockLens_PullFeedsClient) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockLens_PullFeedsClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockLens_PullFeedsClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockLens_PullFeedsClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockLens_PullFeedsClient)(nil).Trailer))
}

// MockLensServer is a mock of LensServer interface.
type MockLensServer struct {
	ctrl     *gomock.Controller
	recorder *MockLensServerMockRecorder
}

// MockLensServerMockRecorder is the mock recorder for MockLensServer.
type MockLensServerMockRecorder struct {
	mock *MockLensServer
}

// NewMockLensServer creates a new mock instance.
func NewMockLensServer(ctrl *gomock.Controller) *MockLensServer {
	mock := &MockLensServer{ctrl: ctrl}
	mock.recorder = &MockLensServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLensServer) EXPECT() *MockLensServerMockRecorder {
	return m.recorder
}

// AddFeed mocks base method.
func (m *MockLensServer) AddFeed(arg0 context.Context, arg1 *api.AddFeedRequest) (*api.AddFeedResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFeed", arg0, arg1)
	ret0, _ := ret[0].(*api.AddFeedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddFeed indicates an expected call of AddFeed.
func (mr *MockLensServerMockRecorder) AddFeed(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFeed", reflect.TypeOf((*MockLensServer)(nil).AddFeed), arg0, arg1)
}

// DeleteFeeds mocks base method.
func (m *MockLensServer) DeleteFeeds(arg0 context.Context, arg1 *api.DeleteFeedsRequest) (*api.DeleteFeedsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFeeds", arg0, arg1)
	ret0, _ := ret[0].(*api.DeleteFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFeeds indicates an expected call of DeleteFeeds.
func (mr *MockLensServerMockRecorder) DeleteFeeds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFeeds", reflect.TypeOf((*MockLensServer)(nil).DeleteFeeds), arg0, arg1)
}

// EditEntries mocks base method.
func (m *MockLensServer) EditEntries(arg0 context.Context, arg1 *api.EditEntriesRequest) (*api.EditEntriesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditEntries", arg0, arg1)
	ret0, _ := ret[0].(*api.EditEntriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditEntries indicates an expected call of EditEntries.
func (mr *MockLensServerMockRecorder) EditEntries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditEntries", reflect.TypeOf((*MockLensServer)(nil).EditEntries), arg0, arg1)
}

// EditFeeds mocks base method.
func (m *MockLensServer) EditFeeds(arg0 context.Context, arg1 *api.EditFeedsRequest) (*api.EditFeedsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditFeeds", arg0, arg1)
	ret0, _ := ret[0].(*api.EditFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditFeeds indicates an expected call of EditFeeds.
func (mr *MockLensServerMockRecorder) EditFeeds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditFeeds", reflect.TypeOf((*MockLensServer)(nil).EditFeeds), arg0, arg1)
}

// ExportOPML mocks base method.
func (m *MockLensServer) ExportOPML(arg0 context.Context, arg1 *api.ExportOPMLRequest) (*api.ExportOPMLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportOPML", arg0, arg1)
	ret0, _ := ret[0].(*api.ExportOPMLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportOPML indicates an expected call of ExportOPML.
func (mr *MockLensServerMockRecorder) ExportOPML(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportOPML", reflect.TypeOf((*MockLensServer)(nil).ExportOPML), arg0, arg1)
}

// GetEntry mocks base method.
func (m *MockLensServer) GetEntry(arg0 context.Context, arg1 *api.GetEntryRequest) (*api.GetEntryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", arg0, arg1)
	ret0, _ := ret[0].(*api.GetEntryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockLensServerMockRecorder) GetEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockLensServer)(nil).GetEntry), arg0, arg1)
}

// GetInfo mocks base method.
func (m *MockLensServer) GetInfo(arg0 context.Context, arg1 *api.GetInfoRequest) (*api.GetInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfo", arg0, arg1)
	ret0, _ := ret[0].(*api.GetInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfo indicates an expected call of GetInfo.
func (mr *MockLensServerMockRecorder) GetInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfo", reflect.TypeOf((*MockLensServer)(nil).GetInfo), arg0, arg1)
}

// GetStats mocks base method.
func (m *MockLensServer) GetStats(arg0 context.Context, arg1 *api.GetStatsRequest) (*api.GetStatsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStats", arg0, arg1)
	ret0, _ := ret[0].(*api.GetStatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStats indicates an expected call of GetStats.
func (mr *MockLensServerMockRecorder) GetStats(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStats", reflect.TypeOf((*MockLensServer)(nil).GetStats), arg0, arg1)
}

// ImportOPML mocks base method.
func (m *MockLensServer) ImportOPML(arg0 context.Context, arg1 *api.ImportOPMLRequest) (*api.ImportOPMLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportOPML", arg0, arg1)
	ret0, _ := ret[0].(*api.ImportOPMLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportOPML indicates an expected call of ImportOPML.
func (mr *MockLensServerMockRecorder) ImportOPML(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportOPML", reflect.TypeOf((*MockLensServer)(nil).ImportOPML), arg0, arg1)
}

// ListEntries mocks base method.
func (m *MockLensServer) ListEntries(arg0 context.Context, arg1 *api.ListEntriesRequest) (*api.ListEntriesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntries", arg0, arg1)
	ret0, _ := ret[0].(*api.ListEntriesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntries indicates an expected call of ListEntries.
func (mr *MockLensServerMockRecorder) ListEntries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntries", reflect.TypeOf((*MockLensServer)(nil).ListEntries), arg0, arg1)
}

// ListFeeds mocks base method.
func (m *MockLensServer) ListFeeds(arg0 context.Context, arg1 *api.ListFeedsRequest) (*api.ListFeedsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeeds", arg0, arg1)
	ret0, _ := ret[0].(*api.ListFeedsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeeds indicates an expected call of ListFeeds.
func (mr *MockLensServerMockRecorder) ListFeeds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeeds", reflect.TypeOf((*MockLensServer)(nil).ListFeeds), arg0, arg1)
}

// PullFeeds mocks base method.
func (m *MockLensServer) PullFeeds(arg0 *api.PullFeedsRequest, arg1 api.Lens_PullFeedsServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullFeeds", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullFeeds indicates an expected call of PullFeeds.
func (mr *MockLensServerMockRecorder) PullFeeds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullFeeds", reflect.TypeOf((*MockLensServer)(nil).PullFeeds), arg0, arg1)
}

// mustEmbedUnimplementedLensServer mocks base method.
func (m *MockLensServer) mustEmbedUnimplementedLensServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedLensServer")
}

// mustEmbedUnimplementedLensServer indicates an expected call of mustEmbedUnimplementedLensServer.
func (mr *MockLensServerMockRecorder) mustEmbedUnimplementedLensServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedLensServer", reflect.TypeOf((*MockLensServer)(nil).mustEmbedUnimplementedLensServer))
}

// MockUnsafeLensServer is a mock of UnsafeLensServer interface.
type MockUnsafeLensServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeLensServerMockRecorder
}

// MockUnsafeLensServerMockRecorder is the mock recorder for MockUnsafeLensServer.
type MockUnsafeLensServerMockRecorder struct {
	mock *MockUnsafeLensServer
}

// NewMockUnsafeLensServer creates a new mock instance.
func NewMockUnsafeLensServer(ctrl *gomock.Controller) *MockUnsafeLensServer {
	mock := &MockUnsafeLensServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeLensServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeLensServer) EXPECT() *MockUnsafeLensServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedLensServer mocks base method.
func (m *MockUnsafeLensServer) mustEmbedUnimplementedLensServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedLensServer")
}

// mustEmbedUnimplementedLensServer indicates an expected call of mustEmbedUnimplementedLensServer.
func (mr *MockUnsafeLensServerMockRecorder) mustEmbedUnimplementedLensServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedLensServer", reflect.TypeOf((*MockUnsafeLensServer)(nil).mustEmbedUnimplementedLensServer))
}

// MockLens_PullFeedsServer is a mock of Lens_PullFeedsServer interface.
type MockLens_PullFeedsServer struct {
	ctrl     *gomock.Controller
	recorder *MockLens_PullFeedsServerMockRecorder
}

// MockLens_PullFeedsServerMockRecorder is the mock recorder for MockLens_PullFeedsServer.
type MockLens_PullFeedsServerMockRecorder struct {
	mock *MockLens_PullFeedsServer
}

// NewMockLens_PullFeedsServer creates a new mock instance.
func NewMockLens_PullFeedsServer(ctrl *gomock.Controller) *MockLens_PullFeedsServer {
	mock := &MockLens_PullFeedsServer{ctrl: ctrl}
	mock.recorder = &MockLens_PullFeedsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLens_PullFeedsServer) EXPECT() *MockLens_PullFeedsServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockLens_PullFeedsServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockLens_PullFeedsServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m_2 *MockLens_PullFeedsServer) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockLens_PullFeedsServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockLens_PullFeedsServer) Send(arg0 *api.PullFeedsResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockLens_PullFeedsServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockLens_PullFeedsServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockLens_PullFeedsServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockLens_PullFeedsServer) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockLens_PullFeedsServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockLens_PullFeedsServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockLens_PullFeedsServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockLens_PullFeedsServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockLens_PullFeedsServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockLens_PullFeedsServer)(nil).SetTrailer), arg0)
}