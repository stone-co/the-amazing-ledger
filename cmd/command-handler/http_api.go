package main

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
)

func httpAPIStart(config configuration.HTTPConfig, log *logrus.Logger, useCase *usecase.LedgerUseCase) {
	accountsHandler := accounts.NewAccountsHandler(log, useCase)
	transactionsHandler := transactions.NewHandler(log, useCase)
	api := http.NewApi(log, accountsHandler, transactionsHandler)
	api.Start("0.0.0.0", config)
}
