// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: iris.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IrisClient is the client API for Iris service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IrisClient interface {
	// AddFeeds adds a new feed source.
	AddFeed(ctx context.Context, in *AddFeedRequest, opts ...grpc.CallOption) (*AddFeedResponse, error)
	// EditFeeds sets one or more fields of feeds.
	EditFeeds(ctx context.Context, in *EditFeedsRequest, opts ...grpc.CallOption) (*EditFeedsResponse, error)
	// ListFeeds lists all added feed sources.
	ListFeeds(ctx context.Context, in *ListFeedsRequest, opts ...grpc.CallOption) (*ListFeedsResponse, error)
	// PullFeeds checks all feeds for updates and returns all unread entries.
	PullFeeds(ctx context.Context, in *PullFeedsRequest, opts ...grpc.CallOption) (Iris_PullFeedsClient, error)
	// DeleteFeeds removes one or more feed sources.
	DeleteFeeds(ctx context.Context, in *DeleteFeedsRequest, opts ...grpc.CallOption) (*DeleteFeedsResponse, error)
	// EditEntries sets one or more fields of an entry.
	EditEntries(ctx context.Context, in *EditEntriesRequest, opts ...grpc.CallOption) (*EditEntriesResponse, error)
	// ExportOPML exports feed subscriptions as an OPML document.
	ExportOPML(ctx context.Context, in *ExportOPMLRequest, opts ...grpc.CallOption) (*ExportOPMLResponse, error)
	// ImportOPML imports an OPML document.
	ImportOPML(ctx context.Context, in *ImportOPMLRequest, opts ...grpc.CallOption) (*ImportOPMLResponse, error)
	// GetInfo returns the version info of the running server.
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
}

type irisClient struct {
	cc grpc.ClientConnInterface
}

func NewIrisClient(cc grpc.ClientConnInterface) IrisClient {
	return &irisClient{cc}
}

func (c *irisClient) AddFeed(ctx context.Context, in *AddFeedRequest, opts ...grpc.CallOption) (*AddFeedResponse, error) {
	out := new(AddFeedResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/AddFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) EditFeeds(ctx context.Context, in *EditFeedsRequest, opts ...grpc.CallOption) (*EditFeedsResponse, error) {
	out := new(EditFeedsResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/EditFeeds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) ListFeeds(ctx context.Context, in *ListFeedsRequest, opts ...grpc.CallOption) (*ListFeedsResponse, error) {
	out := new(ListFeedsResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/ListFeeds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) PullFeeds(ctx context.Context, in *PullFeedsRequest, opts ...grpc.CallOption) (Iris_PullFeedsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Iris_ServiceDesc.Streams[0], "/iris.Iris/PullFeeds", opts...)
	if err != nil {
		return nil, err
	}
	x := &irisPullFeedsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Iris_PullFeedsClient interface {
	Recv() (*PullFeedsResponse, error)
	grpc.ClientStream
}

type irisPullFeedsClient struct {
	grpc.ClientStream
}

func (x *irisPullFeedsClient) Recv() (*PullFeedsResponse, error) {
	m := new(PullFeedsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *irisClient) DeleteFeeds(ctx context.Context, in *DeleteFeedsRequest, opts ...grpc.CallOption) (*DeleteFeedsResponse, error) {
	out := new(DeleteFeedsResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/DeleteFeeds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) EditEntries(ctx context.Context, in *EditEntriesRequest, opts ...grpc.CallOption) (*EditEntriesResponse, error) {
	out := new(EditEntriesResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/EditEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) ExportOPML(ctx context.Context, in *ExportOPMLRequest, opts ...grpc.CallOption) (*ExportOPMLResponse, error) {
	out := new(ExportOPMLResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/ExportOPML", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) ImportOPML(ctx context.Context, in *ImportOPMLRequest, opts ...grpc.CallOption) (*ImportOPMLResponse, error) {
	out := new(ImportOPMLResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/ImportOPML", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *irisClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := c.cc.Invoke(ctx, "/iris.Iris/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IrisServer is the server API for Iris service.
// All implementations must embed UnimplementedIrisServer
// for forward compatibility
type IrisServer interface {
	// AddFeeds adds a new feed source.
	AddFeed(context.Context, *AddFeedRequest) (*AddFeedResponse, error)
	// EditFeeds sets one or more fields of feeds.
	EditFeeds(context.Context, *EditFeedsRequest) (*EditFeedsResponse, error)
	// ListFeeds lists all added feed sources.
	ListFeeds(context.Context, *ListFeedsRequest) (*ListFeedsResponse, error)
	// PullFeeds checks all feeds for updates and returns all unread entries.
	PullFeeds(*PullFeedsRequest, Iris_PullFeedsServer) error
	// DeleteFeeds removes one or more feed sources.
	DeleteFeeds(context.Context, *DeleteFeedsRequest) (*DeleteFeedsResponse, error)
	// EditEntries sets one or more fields of an entry.
	EditEntries(context.Context, *EditEntriesRequest) (*EditEntriesResponse, error)
	// ExportOPML exports feed subscriptions as an OPML document.
	ExportOPML(context.Context, *ExportOPMLRequest) (*ExportOPMLResponse, error)
	// ImportOPML imports an OPML document.
	ImportOPML(context.Context, *ImportOPMLRequest) (*ImportOPMLResponse, error)
	// GetInfo returns the version info of the running server.
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	mustEmbedUnimplementedIrisServer()
}

// UnimplementedIrisServer must be embedded to have forward compatible implementations.
type UnimplementedIrisServer struct {
}

func (UnimplementedIrisServer) AddFeed(context.Context, *AddFeedRequest) (*AddFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFeed not implemented")
}
func (UnimplementedIrisServer) EditFeeds(context.Context, *EditFeedsRequest) (*EditFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditFeeds not implemented")
}
func (UnimplementedIrisServer) ListFeeds(context.Context, *ListFeedsRequest) (*ListFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFeeds not implemented")
}
func (UnimplementedIrisServer) PullFeeds(*PullFeedsRequest, Iris_PullFeedsServer) error {
	return status.Errorf(codes.Unimplemented, "method PullFeeds not implemented")
}
func (UnimplementedIrisServer) DeleteFeeds(context.Context, *DeleteFeedsRequest) (*DeleteFeedsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFeeds not implemented")
}
func (UnimplementedIrisServer) EditEntries(context.Context, *EditEntriesRequest) (*EditEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditEntries not implemented")
}
func (UnimplementedIrisServer) ExportOPML(context.Context, *ExportOPMLRequest) (*ExportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportOPML not implemented")
}
func (UnimplementedIrisServer) ImportOPML(context.Context, *ImportOPMLRequest) (*ImportOPMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportOPML not implemented")
}
func (UnimplementedIrisServer) GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedIrisServer) mustEmbedUnimplementedIrisServer() {}

// UnsafeIrisServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IrisServer will
// result in compilation errors.
type UnsafeIrisServer interface {
	mustEmbedUnimplementedIrisServer()
}

func RegisterIrisServer(s grpc.ServiceRegistrar, srv IrisServer) {
	s.RegisterService(&Iris_ServiceDesc, srv)
}

func _Iris_AddFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).AddFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/AddFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).AddFeed(ctx, req.(*AddFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_EditFeeds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditFeedsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).EditFeeds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/EditFeeds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).EditFeeds(ctx, req.(*EditFeedsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_ListFeeds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFeedsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).ListFeeds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/ListFeeds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).ListFeeds(ctx, req.(*ListFeedsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_PullFeeds_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PullFeedsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IrisServer).PullFeeds(m, &irisPullFeedsServer{stream})
}

type Iris_PullFeedsServer interface {
	Send(*PullFeedsResponse) error
	grpc.ServerStream
}

type irisPullFeedsServer struct {
	grpc.ServerStream
}

func (x *irisPullFeedsServer) Send(m *PullFeedsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Iris_DeleteFeeds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFeedsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).DeleteFeeds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/DeleteFeeds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).DeleteFeeds(ctx, req.(*DeleteFeedsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_EditEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).EditEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/EditEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).EditEntries(ctx, req.(*EditEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_ExportOPML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportOPMLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).ExportOPML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/ExportOPML",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).ExportOPML(ctx, req.(*ExportOPMLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_ImportOPML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportOPMLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).ImportOPML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/ImportOPML",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).ImportOPML(ctx, req.(*ImportOPMLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Iris_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IrisServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iris.Iris/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IrisServer).GetInfo(ctx, req.(*GetInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Iris_ServiceDesc is the grpc.ServiceDesc for Iris service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Iris_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iris.Iris",
	HandlerType: (*IrisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddFeed",
			Handler:    _Iris_AddFeed_Handler,
		},
		{
			MethodName: "EditFeeds",
			Handler:    _Iris_EditFeeds_Handler,
		},
		{
			MethodName: "ListFeeds",
			Handler:    _Iris_ListFeeds_Handler,
		},
		{
			MethodName: "DeleteFeeds",
			Handler:    _Iris_DeleteFeeds_Handler,
		},
		{
			MethodName: "EditEntries",
			Handler:    _Iris_EditEntries_Handler,
		},
		{
			MethodName: "ExportOPML",
			Handler:    _Iris_ExportOPML_Handler,
		},
		{
			MethodName: "ImportOPML",
			Handler:    _Iris_ImportOPML_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _Iris_GetInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PullFeeds",
			Handler:       _Iris_PullFeeds_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "iris.proto",
}