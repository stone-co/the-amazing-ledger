package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/healthcheck"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
	"github.com/urfave/negroni"
)

type Api struct {
	log          *logrus.Logger
	Healthcheck  healthcheck.Handler
	Accounts     *accounts.Handler
	Transactions *transactions.Handler
	//	Middleware  common.Middleware
}

func NewApi(log *logrus.Logger, accounts *accounts.Handler, transactions *transactions.Handler) *Api {
	return &Api{
		log:          log,
		Accounts:     accounts,
		Transactions: transactions,
	}
}

func (a *Api) Start(host string, cfg configuration.HTTPConfig) {
	// Router
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/health", a.Healthcheck.Get).Methods("GET")
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")

	//Accounts
	r.HandleFunc("/accounts", a.Accounts.Create).Methods("POST")

	//Entries
	r.HandleFunc("/transactions", a.Transactions.Create).Methods("POST")

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(r)

	endpoint := fmt.Sprintf("%s:%d", host, cfg.Port)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Handler: n,
		Addr:    endpoint,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		a.log.Infof("starting http api at %s", endpoint)
		serverErrors <- srv.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		a.log.WithError(err).Fatal("http server error")

	case sig := <-shutdown:
		a.log.Printf("%v : Start http api shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			a.log.WithError(err).Fatal("could not stop http server gracefully")
		}
		a.log.Printf("%v : Finished http api shutdown", sig)
	}
}
