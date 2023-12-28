// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"encoding/xml"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bow/neon/internal/entity"
)

func errorUnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	rsp, err := handler(ctx, req)
	return rsp, mapError(err)
}

func errorStreamServerInterceptor(
	srv any,
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	return mapError(handler(srv, ss))
}

func unwrapError(err error) (codes.Code, error) {
	if err == nil {
		return codes.Unknown, nil
	}
	switch cerr := err.(type) {
	case entity.FeedNotFoundError, entity.EntryNotFoundError:
		return codes.NotFound, cerr
	case xml.UnmarshalError, *xml.SyntaxError:
		return codes.InvalidArgument, cerr
	default:
		var (
			ierr  error
			icode codes.Code
		)
		if uerr := errors.Unwrap(err); uerr != nil {
			icode, ierr = unwrapError(uerr)
		}
		if ierr != nil {
			return icode, ierr
		}
		return codes.Unknown, err
	}
}

func mapError(err error) error {
	code, err := unwrapError(err)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.Unknown {
			return st.Err()
		}
		return status.Error(code, err.Error())
	}
	return nil
}
