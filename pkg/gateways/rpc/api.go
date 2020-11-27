package rpc

import (
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.LedgerServiceServer = &API{}

type API struct {
	log     *logrus.Logger
	UseCase ledger.UseCase
}

func NewAPI(log *logrus.Logger, useCase ledger.UseCase) *API {
	return &API{
		log:     log,
		UseCase: useCase,
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

	proto.RegisterLedgerServiceServer(srv, a)
	proto.RegisterHealthServer(srv, a)

	return srv
}
