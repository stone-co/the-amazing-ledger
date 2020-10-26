package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions"
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

func (a *API) Start(cfg configuration.GRPCConfig) {
	endpoint := fmt.Sprintf(":%d", cfg.Port)

	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a.log.Infof("starting grpc api at %s", endpoint)
	server := grpc.NewServer()

	proto.RegisterLedgerServiceServer(server, a.handler)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
