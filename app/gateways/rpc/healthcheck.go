package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (API) Check(_ context.Context, _ *empty.Empty) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING_STATUS_SERVING,
	}, nil
}
