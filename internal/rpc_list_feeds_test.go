package internal

import (
	"context"
	"testing"

	"github.com/bow/courier/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFeedsOkEmpty(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	req := api.ListFeedsRequest{}

	client := newTestClientBuilder(t).Build()

	rsp, err := client.ListFeeds(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Empty(rsp.GetFeeds())
}
