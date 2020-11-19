package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	httpApi "github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
)

func NewHttpServer(config configuration.HTTPConfig, log *logrus.Logger, useCase *usecase.LedgerUseCase) *http.Server {
	transactionsHandler := transactions.NewHandler(log, useCase)
	api := httpApi.NewApi(log, transactionsHandler)
	return api.NewServer("0.0.0.0", config)
}
