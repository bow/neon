package internal

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bow/courier/api"
	gomock "github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func TestAddFeedOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	parser := NewMockFeedParser(gomock.NewController(t))
	parser.
		EXPECT().
		ParseURL("https://bar.com/feed.xml").
		MaxTimes(1).
		Return(&gofeed.Feed{}, nil)

	storePath := filepath.Join(t.TempDir(), "courier-add-feed.db")
	r.NoFileExists(storePath)

	server := defaultTestServerBuilder(t).Parser(parser).StorePath(storePath)
	client := newTestClientBuilder().ServerBuilder(server).Build(t)
	r.FileExists(storePath)

	req := api.AddFeedRequest{
		Url:        "https://bar.com/feed.xml",
		Categories: []string{"c1", "c2"},
	}
	rsp, err := client.AddFeed(context.Background(), &req)

	r.NoError(err)
	r.NotNil(rsp)
}
