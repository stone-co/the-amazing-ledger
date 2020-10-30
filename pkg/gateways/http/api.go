package http

import (
	"fmt"
	"net/http"

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

func (a *Api) NewServer(host string, cfg configuration.HTTPConfig) *http.Server {
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

	srv := &http.Server{
		Handler: n,
		Addr:    endpoint,
	}

	return srv
}
