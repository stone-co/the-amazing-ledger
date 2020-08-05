package accounts

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts"
)

type Handler struct {
	log     *logrus.Logger
	UseCase accounts.UseCase
}

func NewAccountsHandler(log *logrus.Logger, useCase accounts.UseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
