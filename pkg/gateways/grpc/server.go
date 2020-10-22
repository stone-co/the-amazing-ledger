package grpc

import (
	"log"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions"
)

type Server struct {
	log     *logrus.Logger
	handler *transactions.Handler
}

func NewServer(log *logrus.Logger, handler *transactions.Handler) *Server {
	return &Server{
		log:     log,
		handler: handler,
	}
}

func (s *Server) Start(cfg configuration.GRPCConfig) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.log.Infof("starting grpc server at %d port", cfg.Port)
	server := grpc.NewServer()

	pb.RegisterLedgerServiceServer(server, s.handler)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
