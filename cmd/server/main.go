package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

func main() {
	log := logrus.New()
	log.Infoln("starting ledger process...")
	log.Infof("build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)

	cfg, err := app.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to load app configuration")
	}

	nr, err := newrelic.App(cfg.NewRelic.AppName, cfg.NewRelic.LicenseKey, logrus.NewEntry(log))
	if err != nil {
		log.WithError(err).Fatal("error starting new relic")
	}

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), log)
	if err != nil {
		log.WithError(err).Fatal("unable to connect to database")
	}
	defer conn.Close()

	if err = postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		log.WithError(err).Panic("running postgres migrations")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port))
	if err != nil {
		log.WithError(err).Panic("failed to listen")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, log)
	ledgerUseCase := usecases.NewLedgerUseCase(log, ledgerRepository)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	rpcServer, gwServer, err := rpc.NewServer(ctx, ledgerUseCase, nr, cfg, log, BuildGitCommit, BuildTime)
	if err != nil {
		log.WithError(err).Panic("failed to create server")
	}

	go func() {
		if err = rpcServer.Serve(lis); err != nil {
			log.WithError(err).Error("rpcServer Serve failed")
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
		defer cancel()

		rpcServer.GracefulStop()

		if err = gwServer.Shutdown(ctx); err != nil {
			_ = gwServer.Close()
			log.WithError(err).Fatal("could not stop server gracefully")
		}
	}()

	go handleInterrupt(log, cancel)

	err = gwServer.ListenAndServe()
	if err != nil {
		log.WithError(err).Error("gwServer ListenAndServe failed")
	}
}

func handleInterrupt(log *logrus.Logger, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	sig := <-signals
	log.Infof("captured signal: %v - server shutdown\n", sig)
	signal.Stop(signals)
	cancel()
}
