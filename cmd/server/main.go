package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/probes"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	server "github.com/stone-co/the-amazing-ledger/app/gateways/http"
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

	newrelic, err := newrelic.NewNewrelic(cfg.NewRelic.AppName, cfg.NewRelic.LicenseKey, logrus.NewEntry(log))
	if err != nil {
		log.WithError(err).Fatal("error starting new relic")
	}

	ledgerProbe := probes.NewLedgerProbe(log, newrelic)

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), log)
	if err != nil {
		log.WithError(err).Fatal("unable to connect to database")
	}
	defer conn.Close()

	if err = postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		log.WithError(err).Fatal("running postgres migrations")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, ledgerProbe)
	ledgerUseCase := usecases.NewLedgerUseCase(ledgerRepository, ledgerProbe)

	httpServer := server.NewHttpServer(cfg.HttpServer, BuildGitCommit, BuildTime, log)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
		defer cancel()
		if er := httpServer.Shutdown(ctx); er != nil {
			_ = httpServer.Close()
			log.WithError(er).Fatal("server error:could not stop server gracefully")
		}
	}()
	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	// Initialize the server (grpc-gateway)
	rpcServer, err := NewGRPCServer(ledgerUseCase, newrelic, cfg.RPCServer, log)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize the server")
	}
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	// NewServer Server listening for requests.
	go func() {
		log.Infof("ready to accept connections at: %s", rpcServer.Addr)
		serverErrors <- fmt.Errorf("server's ListenAndServe failed. %w", rpcServer.ListenAndServe())
	}()

	// =================
	// Shutdown

	//Handle OS signals
	go handleInterrupt(cfg, log, rpcServer)

	// Blocking main and waiting for server error.
	err = <-serverErrors
	if !errors.Is(err, http.ErrServerClosed) {
		log.WithError(fmt.Errorf("server error: %w", err)).Fatal()
	}
}

func handleInterrupt(cfg *app.Config, log *logrus.Logger, sv *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	sig := <-signals
	log.Infof("captured signal: %v - server shutdown\n", sig)
	signal.Stop(signals)
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RPCServer.ShutdownTimeout)
	defer cancel()
	// Asking listener to shutdown and shed load.
	if err := sv.Shutdown(ctx); err != nil {
		_ = sv.Close()
		log.WithError(err).Fatal("server error:could not stop server gracefully")
	}
}
