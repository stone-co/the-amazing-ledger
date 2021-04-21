package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app"
)

func NewHttpServer(s app.HttpServerConfig, commit string, time string, log *logrus.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version(w, commit, time, log)
	}))

	log.WithFields(logrus.Fields{
		"port": strconv.Itoa(s.Port),
	}).Infoln("starting http server...")

	return &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.Port),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
}

func version(w http.ResponseWriter, commit string, time string, log *logrus.Logger) {
	b, err := json.Marshal(struct {
		GitCommitHash string `json:"git_hash"`
		BUildTime     string `json:"time"`
	}{
		GitCommitHash: commit,
		BUildTime:     time,
	})
	if err != nil {
		log.WithError(err).Errorln("Failed to marshal version")
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(b)
	if err != nil {
		log.WithError(err).Errorln("Failed to write version body")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
