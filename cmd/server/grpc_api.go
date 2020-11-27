package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func NewGRPCServer(useCase *usecase.LedgerUseCase, cfg configuration.ServerConfig, log *logrus.Logger) (*http.Server, error) {
	api := rpc.NewAPI(log, useCase, useCase)
	grpcServer := api.NewServer()
	server, err := NewServer(grpcServer, cfg)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func NewServer(grpcServer *grpc.Server, cfg configuration.ServerConfig) (*http.Server, error) {
	// gwMux is the grpc-gateway ServeMux, used to serve HTTP/REST requests.
	gwMux := runtime.NewServeMux()
	// Always use localhost for gateway
	gwEndpoint := fmt.Sprintf("localhost:%d", cfg.Port)

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
