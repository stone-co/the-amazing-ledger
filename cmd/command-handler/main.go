package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	accountsUsecase "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/usecase"
	transactionsUsecase "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
)

func main() {
	log := logrus.New()
	log.Infoln("Starting Command-Handler process...")

	cfg, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load app configuration: %s", err.Error()))
	}

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), log)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	defer conn.Close()

	if err := postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		log.WithError(err).Fatal("Error running postgres migrations")
	}

	accountsRepository := postgres.NewAccountsRepository(conn, log)
	accountsUseCase := accountsUsecase.NewAccountsUseCase(log, accountsRepository)
	accountsHandler := accounts.NewAccountsHandler(log, accountsUseCase)

	transactionsRepository := postgres.NewTransactionsRepository(conn, log)
	transactionsUseCase := transactionsUsecase.NewTransactionsUseCase(log, transactionsRepository)
	transactionsHandler := transactions.NewTransactionsHandler(log, transactionsUseCase)

	// Starting gateway http API
	api := http.NewApi(log, accountsHandler, transactionsHandler)
	api.Start("0.0.0.0", cfg.API.Port)
}
