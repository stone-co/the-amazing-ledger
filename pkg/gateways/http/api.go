package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/healthcheck"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/transactions"
	"github.com/urfave/negroni"
)

type Api struct {
	log          *logrus.Logger
	Healthcheck  healthcheck.Handler
	Transactions *transactions.Handler
}

func NewApi(log *logrus.Logger, transactions *transactions.Handler) *Api {
	return &Api{
		log:          log,
		Transactions: transactions,
	}
}

func (a *Api) NewServer(host string, cfg configuration.HTTPConfig) *http.Server {
	// Router
	r := mux.NewRouter().UseEncodedPath()

	// Handlers
	r.HandleFunc("/health", a.Healthcheck.Get).Methods("GET")
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods("GET")

	//Accounts
	r.HandleFunc("/accounts/{account_path}/balance", a.Transactions.GetAccountBalance).Methods("GET")

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
