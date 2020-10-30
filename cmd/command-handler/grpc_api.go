package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	rpc "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions"
)

func NewGRPCServer(log *logrus.Logger, useCase *usecase.LedgerUseCase) *grpc.Server {
	transactionsHandler := transactions.NewHandler(log, useCase)
	api := rpc.NewAPI(log, transactionsHandler)
	return api.NewServer()
}
