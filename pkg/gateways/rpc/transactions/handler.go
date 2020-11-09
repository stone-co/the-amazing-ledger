package transactions

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/proto"
)

var _ proto.LedgerServiceServer = &Handler{}

type Handler struct {
	log     *logrus.Logger
	UseCase ledger.TransactionsUseCase
}

func NewHandler(log *logrus.Logger, useCase ledger.TransactionsUseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
