package internal

import (
	"context"
	"testing"

	"github.com/bow/courier/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfoOk(t *testing.T) {
	client := setupTestServer(t)

	req := api.GetInfoRequest{}
	rsp, err := client.GetInfo(context.Background(), &req)
	require.NoError(t, err)

	want := &api.GetInfoResponse{
		Name:      AppName(),
		Version:   Version(),
		GitCommit: GitCommit(),
		BuildTime: BuildTime(),
	}
	a := assert.New(t)
	a.Equal(want.Name, rsp.Name)
	a.Equal(want.Version, rsp.Version)
	a.Equal(want.GitCommit, rsp.GitCommit)
	a.Equal(want.BuildTime, rsp.BuildTime)
}
