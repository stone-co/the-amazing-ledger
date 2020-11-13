package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
)

func NewInternal(s configuration.PrometheusServerConfig) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.Port),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
}
