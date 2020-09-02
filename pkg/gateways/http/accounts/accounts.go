package accounts

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
)

type Handler struct {
	log     *logrus.Logger
	UseCase ledger.AccountsUseCase
}

func NewAccountsHandler(log *logrus.Logger, useCase ledger.AccountsUseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
