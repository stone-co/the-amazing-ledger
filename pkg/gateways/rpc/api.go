package rpc

import (
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/proto"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/transactions"
)

type API struct {
	log     *logrus.Logger
	handler *transactions.Handler
}

func NewAPI(log *logrus.Logger, handler *transactions.Handler) *API {
	return &API{
		log:     log,
		handler: handler,
	}
}

func (a *API) NewServer() *grpc.Server {
	// Define a func to handle panic
	dealPanic := func(p interface{}) (err error) {
		log.Printf("panic triggered: %v", p)
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(dealPanic),
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

	proto.RegisterLedgerServiceServer(srv, a.handler)

	return srv
}
