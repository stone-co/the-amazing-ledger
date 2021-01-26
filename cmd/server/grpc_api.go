package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func NewGRPCServer(useCase *usecases.LedgerUseCase, nr *newrelic.Application, cfg app.ServerConfig, log *logrus.Logger) (*http.Server, error) {
	api := rpc.NewAPI(log, useCase)
	grpcServer := api.NewServer(nr)
	server, err := NewServer(grpcServer, cfg)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func NewServer(grpcServer *grpc.Server, cfg app.ServerConfig) (*http.Server, error) {
	// gwMux is the grpc-gateway ServeMux, used to serve HTTP/REST requests.
	gwMux := runtime.NewServeMux()
	// Always use localhost for gateway
	gwEndpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := proto.RegisterLedgerServiceHandlerFromEndpoint(context.Background(), gwMux, gwEndpoint, opts); err != nil {
		return nil, err
	}
	if err := proto.RegisterHealthHandlerFromEndpoint(context.Background(), gwMux, gwEndpoint, opts); err != nil {
		return nil, err
	}

	// use a single server for both operations
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L55
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
			return
		}
		// add http method info to gRPC gateway server
		r.Header.Add("Grpc-Metadata-HTTP-Method", r.Method)
		gwMux.ServeHTTP(w, r)
	})

	return &http.Server{
		Addr:         gwEndpoint,
		Handler:      h2c.NewHandler(rootHandler, &http2.Server{}),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}, nil
}
