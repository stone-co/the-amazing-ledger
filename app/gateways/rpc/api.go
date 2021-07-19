package rpc

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

var _ proto.LedgerServiceServer = &API{}

type API struct {
	log     *logrus.Logger
	UseCase domain.UseCase
}

func NewAPI(log *logrus.Logger, useCase domain.UseCase) *API {
	return &API{
		log:     log,
		UseCase: useCase,
	}
}
