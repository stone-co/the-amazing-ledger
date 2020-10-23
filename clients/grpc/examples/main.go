package main

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func main() {
	logrus := logrus.New()
	log := logrus.WithField("ClientSDK", "Test")

	log.Println("GRPC example starting...")
	defer log.Println("GRPC example finishing...")

	// Connect to the Ledger gRPC server
	host := "localhost"
	port := 50051
	conn, err := ledger.Connect(host, port)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	defineTransactionWithThreeEntries(log, conn)
	invalidTransactionsTests(log, conn)
}
