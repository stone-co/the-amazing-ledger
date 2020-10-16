package transactions

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions/proto"
)

// TODO: move to outside
type Handler struct {
	log     *logrus.Logger
	UseCase ledger.TransactionsUseCase
}

// TODO: move to outside
func NewHandler(log *logrus.Logger, useCase ledger.TransactionsUseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}

func (s *Handler) SayHello(ctx context.Context, in pb.CreateTransactionInput) (*pb.CreateTransactionResult, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.CreateTransactionResult{Message: "Transaction " + in.GetId()}, nil
}
