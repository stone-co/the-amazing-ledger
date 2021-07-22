package server

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	promhttp.Handler().ServeHTTP(w, r)
}

func VersionHandler(commit, time string) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		b, err := json.Marshal(struct {
			GitCommitHash string `json:"git_hash"`
			BuildTime     string `json:"time"`
		}{
			GitCommitHash: commit,
			BuildTime:     time,
		})
		if err != nil {
			logrus.Errorf("failed to marshal version: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(b)
		if err != nil {
			logrus.Errorf("failed to write version body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
