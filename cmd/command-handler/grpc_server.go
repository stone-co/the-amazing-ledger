package main

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/usecase"
	"github.com/stone-co/the-amazing-ledger/pkg/common/configuration"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/transactions"
)

func grpcServerStart(config configuration.GrpcConfig, log *logrus.Logger, useCase *usecase.LedgerUseCase) {
	transactionsGrpcHandler := transactions.NewHandler(log, useCase)
	grpcServer := grpc.NewServer(log, transactionsGrpcHandler)
	go grpcServer.Start(config)
}
