package grpc

import (
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions/proto"
)

type Server struct {
	log *logrus.Logger
	pb.UnimplementedTransactionsServer
}

func NewServer(log *logrus.Logger) *Server {
	return &Server{
		log: log,
	}
}

func (a *Server) Start(cfg configuration.GrpcConfig) {

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterTransactionsServer(server, &Server{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
