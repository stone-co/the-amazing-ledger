package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/prometheus"
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
		log.WithError(err).Fatal("running postgres migrations")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, log)

	ledgerUseCase := usecase.NewLedgerUseCase(log, ledgerRepository)
	if err = ledgerUseCase.LoadObjectsIntoCache(context.Background()); err != nil {
		log.WithError(err).Fatal("failed to populate cache")
	}

	metricsServer := prometheus.NewInternal(cfg.Metrics.Prometheus)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Metrics.Prometheus.ShutdownTimeout)
		defer cancel()
		if er := metricsServer.Shutdown(ctx); er != nil {
			_ = metricsServer.Close()
			log.WithError(er).Fatal("server error:could not stop server gracefully")
		}
	}()

	// Initialize the server (grpc-gateway)
	server, err := NewGRPCServer(ledgerUseCase, cfg.Server, log)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize the server")
	}
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	// NewServer Server listening for requests.
	go func() {
		log.Infof("ready to accept connections at: %s", server.Addr)
		serverErrors <- fmt.Errorf("server's ListenAndServe failed. %w", server.ListenAndServe())
	}()

	// =================
	// Shutdown

	//Handle OS signals
	go handleInterrupt(cfg, log, server)

	// Blocking main and waiting for server error.
	err = <-serverErrors
	if !errors.Is(err, http.ErrServerClosed) {
		log.WithError(fmt.Errorf("server error: %w", err)).Fatal()
	}
}

func handleInterrupt(cfg *configuration.Config, log *logrus.Logger, sv *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	sig := <-signals
	log.Infof("captured signal: %v - server shutdown\n", sig)
	signal.Stop(signals)
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	// Asking listener to shutdown and shed load.
	if err := sv.Shutdown(ctx); err != nil {
		_ = sv.Close()
		log.WithError(err).Fatal("server error:could not stop server gracefully")
	}
}
