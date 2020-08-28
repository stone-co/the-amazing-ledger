package http

import (
	"fmt"
	"net/http"

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

func (a *Api) Start(host, port string) {
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

	endpoint := fmt.Sprintf("%s:%s", host, port)

	srv := &http.Server{
		Handler: n,
		Addr:    endpoint,
	}

	a.log.Infof("starting API at %s", endpoint)

	if err := srv.ListenAndServe(); err != nil {
		a.log.Fatalf("can't run server: %s", err.Error())
	}
}
