package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	httpApi "github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
)

func NewHttpServer(config configuration.HTTPConfig, log *logrus.Logger, useCase *usecase.LedgerUseCase) *http.Server {
	accountsHandler := accounts.NewAccountsHandler(log, useCase)
	transactionsHandler := transactions.NewHandler(log, useCase)
	api := httpApi.NewApi(log, accountsHandler, transactionsHandler)
	return api.NewServer("0.0.0.0", config)
}
