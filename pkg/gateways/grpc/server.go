package grpc

import (
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
)

type Server struct {
	log *logrus.Logger
	pb.UnimplementedLedgerServiceServer
}

func NewServer(log *logrus.Logger) *Server {
	return &Server{
		log: log,
	}
}

func (a *Server) Start(cfg configuration.GrpcConfig) {

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a.log.Infof("starting grpc server at %s port", cfg.Port)
	server := grpc.NewServer()

	pb.RegisterLedgerServiceServer(server, &Server{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
