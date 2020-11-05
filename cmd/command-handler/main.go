package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/db/postgres"
	"golang.org/x/sync/errgroup"
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
	httpServer := NewHttpServer(cfg.API, log, ledgerUseCase)
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
		log.Fatal(fmt.Errorf("server error: %w", err))

	case sig := <-shutdown:
		log.Printf("NewServer shutdown %v\n", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		group := &errgroup.Group{}

		group.Go(func() error {
			log.Printf("Stopping HTTP Server %v\n", sig)
			defer log.Printf("HTTP Server Stopped %v\n", sig)
			// Asking listener to shutdown and shed load.
			if err := httpServer.Shutdown(ctx); err != nil {
				_ = httpServer.Close()
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
			return nil
		})
		group.Go(func() error {
			log.Printf("Stopping GRCP Server %v\n", sig)
			defer log.Printf("GRCP Server Stopped %v\n", sig)
			stopped := make(chan struct{})
			go func() {
				grpcServer.GracefulStop()
				close(stopped)
			}()

			t := time.NewTimer(cfg.ShutdownTimeout)
			select {
			case <-t.C:
				grpcServer.Stop()
				return errors.New("could not stop grpc server gracefully")
			case <-stopped:
				t.Stop()
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			log.Fatal(fmt.Errorf("server error: %w", err))
		}
	}
}
