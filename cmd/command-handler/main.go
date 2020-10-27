package main

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		grpcAPIStart(cfg.GRPC, log, ledgerUseCase)
	}()

	go func() {
		defer wg.Done()
		httpAPIStart(cfg.API, log, ledgerUseCase)
	}()

	wg.Wait()
}
