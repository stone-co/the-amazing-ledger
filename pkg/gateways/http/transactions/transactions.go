package transactions

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions"
)

type Handler struct {
	log     *logrus.Logger
	UseCase transactions.UseCase
}

func NewTransactionsHandler(log *logrus.Logger, useCase transactions.UseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
