package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions"
)

type API struct {
	log     *logrus.Logger
	handler *transactions.Handler
}

func NewAPI(log *logrus.Logger, handler *transactions.Handler) *API {
	return &API{
		log:     log,
		handler: handler,
	}
}

func (a *API) Start(cfg configuration.GRPCConfig) {
	endpoint := fmt.Sprintf(":%d", cfg.Port)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterLedgerServiceServer(srv, a.handler)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		a.log.Infof("starting grpc api at %s", endpoint)
		serverErrors <- srv.Serve(lis)
	}()

	// =========================================================================
	// Shutdown

	// Blocking and waiting for shutdown.
	select {
	case err := <-serverErrors:
		a.log.WithError(err).Fatal("grpc server error")

	case sig := <-shutdown:
		a.log.Printf("%v : Start grpc api shutdown", sig)
		a.gracefulStopWithTimeout(srv, cfg.ShutdownTimeout)
		a.log.Printf("%v : Finished grpc api shutdown", sig)
	}
}

func (a API) gracefulStopWithTimeout(srv *grpc.Server, timeout time.Duration) {
	stopped := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(stopped)
	}()

	t := time.NewTimer(timeout)
	select {
	case <-t.C:
		srv.Stop()
		a.log.Fatal("could not stop grpc server gracefully")
	case <-stopped:
		t.Stop()
	}
}
