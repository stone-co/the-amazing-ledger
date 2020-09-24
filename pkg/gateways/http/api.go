package http

import (
	"context"
	"fmt"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/entries"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/healthcheck"
	"github.com/urfave/negroni"
)

type Api struct {
	log         *logrus.Logger
	Healthcheck healthcheck.Handler
	Accounts    *accounts.Handler
	Entries     *entries.Handler
	//	Middleware  common.Middleware
}

func NewApi(log *logrus.Logger, accounts *accounts.Handler, entries *entries.Handler) *Api {
	return &Api{
		log:      log,
		Accounts: accounts,
		Entries:  entries,
	}
}

func (a *Api) Start(host string, cfg configuration.APIConfig) {
	// Router
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/health", a.Healthcheck.Get).Methods("GET")
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")

	//Accounts
	r.HandleFunc("/accounts", a.Accounts.Create).Methods("POST")

	//Entries
	r.HandleFunc("/transactions", a.Entries.Create).Methods("POST")

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(r)

	endpoint := fmt.Sprintf("%s:%s", host, cfg.Port)

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
		a.log.Infof("starting API at %s", endpoint)
		serverErrors <- srv.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		a.log.WithError(err).Fatal("server error")

	case sig := <-shutdown:
		a.log.Printf("%v : Start shutdown", sig)
		timeout, err := time.ParseDuration(cfg.ShutdownTimeout)
		if err != nil {
			a.log.Warnf("error parsing duration %q. fallbacking to 5 seconds\n", cfg.ShutdownTimeout)
			timeout = 5 * time.Second
		}
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			a.log.WithError(err).Fatal("could not stop server gracefully")
		}
	}
}
