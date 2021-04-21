package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app"
)

func NewHttpServer(s app.HttpServerConfig, commit string, time string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version(w, commit, time)
	}))

	log := logrus.New()
	log.WithFields(logrus.Fields{
		"port": strconv.Itoa(s.Port),
	}).Infoln("starting ledger process...")

	return &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.Port),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
}

func version(w http.ResponseWriter, commit string, time string) {
	w.Write([]byte(commit + "-" + time))
}
