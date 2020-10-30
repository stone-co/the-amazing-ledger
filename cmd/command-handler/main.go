package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
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
		log.WithError(err).Fatal("running postgres migrations")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, log)

	ledgerUseCase := usecase.NewLedgerUseCase(log, ledgerRepository)
	if err = ledgerUseCase.LoadObjectsIntoCache(context.Background()); err != nil {
		log.WithError(err).Fatal("failed to populate cache")
	}

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// NewServer HTTP Server listening for requests.
	httpServer := httpAPIStart(cfg.API, log, ledgerUseCase)
	go func() {
		log.Infof("starting http api at %s", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// NewServer GRPC Server listening for requests.
	grpcServer := NewGRPCServer(log, ledgerUseCase)
	go func() {
		endpoint := fmt.Sprintf(":%d", cfg.GRPC.Port)
		log.Infof("starting grpc api at %s", endpoint)
		lis, err := net.Listen("tcp", endpoint)
		if err != nil {
			serverErrors <- err
		}
		serverErrors <- grpcServer.Serve(lis)
	}()

	// =================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatal(errors.Wrap(err, "server error"))

	case sig := <-shutdown:
		log.Printf("NewServer shutdown %v\n", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		log.Printf("Stopping HTTP Server %v\n", sig)
		if err := httpServer.Shutdown(ctx); err != nil {
			_ = httpServer.Close()
			log.Fatal(errors.Wrap(err, "could not stop server gracefully"))
		}
		log.Printf("HTTP Server Stopped %v\n", sig)

		log.Printf("Stopping GRCP Server %v\n", sig)
		grpcServer.GracefulStop()
		log.Printf("GRCP Server Stopped %v\n", sig)
	}
}
