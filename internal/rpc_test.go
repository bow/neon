package internal

import (
	"context"
	"testing"

	"github.com/bow/courier/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEditFeedOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.EditFeedRequest{}
	rsp, err := client.EditFeed(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestListFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.ListFeedsRequest{}
	rsp, err := client.ListFeeds(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestDeleteFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.DeleteFeedsRequest{}
	rsp, err := client.DeleteFeeds(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestPollFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	stream, err := client.PollFeeds(context.Background())
	r.NoError(err)
	waitc := make(chan struct{})

	go func() {
		for {
			rsp, errStream := stream.Recv()
			r.Nil(rsp)
			r.EqualError(errStream, status.New(codes.Unimplemented, "unimplemented").String())
			close(waitc)
			return
		}
	}()

	req := api.PollFeedsRequest{}
	err = stream.Send(&req)
	r.NoError(err)

	err = stream.CloseSend()
	r.NoError(err)
	<-waitc
}

func TestEditEntryOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.EditEntryRequest{}
	rsp, err := client.EditEntry(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestExportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.ExportOPMLRequest{}
	rsp, err := client.ExportOPML(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestImportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.ImportOPMLRequest{}
	rsp, err := client.ImportOPML(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}
