package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/courier/internal/store"
)

func storeErrorUnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	rsp, err := handler(ctx, req)
	return rsp, mapStoreError(err)
}

func storeErrorStreamServerInterceptor(
	srv any,
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	return mapStoreError(handler(srv, ss))
}

func mapStoreError(err error) error {
	switch cerr := err.(type) {
	case store.FeedNotFoundError:
		return status.Error(codes.NotFound, cerr.Error())
	case store.EntryNotFoundError:
		return status.Error(codes.NotFound, cerr.Error())
	default:
		return err
	}
}
