package rpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (API) Check(_ context.Context, _ *emptypb.Empty) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING_STATUS_SERVING,
	}, nil
}
