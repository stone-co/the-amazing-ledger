package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
)

func main() {
	log := logrus.New()
	log.Infoln("Starting Ledger process...")

	cfg, err := configuration.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to load app configuration")
	}

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), log)
	if err != nil {
		log.WithError(err).Fatal("unable to connect to database")
	}
	defer conn.Close()

	if err = postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		log.WithError(err).Fatal("error running postgres migrations")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, log)

	ledgerUseCase := usecase.NewLedgerUseCase(log, ledgerRepository)
	if err = ledgerUseCase.LoadObjectsIntoCache(context.Background()); err != nil {
		log.WithError(err).Fatal("failed to populate cache")
	}

	accountsHandler := accounts.NewAccountsHandler(log, ledgerUseCase)
	transactionsHandler := transactions.NewHandler(log, ledgerUseCase)

	// Starting gateway http API
	api := http.NewApi(log, accountsHandler, transactionsHandler)
	api.Start("0.0.0.0", cfg.API)

	grpcServer := grpc.NewServer(log)
	grpcServer.Start(cfg.Grpc)
}
