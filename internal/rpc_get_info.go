package internal

import (
	"context"

	"github.com/bow/courier/api"
)

// GetInfo satisfies the service API.
func (r *rpc) GetInfo(
	_ context.Context,
	_ *api.GetInfoRequest,
) (*api.GetInfoResponse, error) {

	rsp := api.GetInfoResponse{
		Name:      AppName(),
		Version:   Version(),
		GitCommit: GitCommit(),
		BuildTime: BuildTime(),
	}

	return &rsp, nil
}
