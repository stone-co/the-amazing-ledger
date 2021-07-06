package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func main() {
	logrus := logrus.New()
	log := logrus.WithField("ClientSDK", "Test")

	log.Println("Server example starting...")
	defer log.Println("Server example finishing...")

	// Connect to the Ledger gRPC server
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	host := "localhost"
	port := 3000
	conn, err := ledger.Connect(ctx, host, port)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// defineTransactionWithThreeEntries(log, conn)
	// invalidTransactionsTests(log, conn)

	// getAccountBalance(log, conn)
	// getAccountBalanceWithMoreEntries(log, conn)
	// getAccountBalanceNotFoundAccount(log, conn)
	// getAccountBalanceWithWildcard(log, conn)
	// getAnalyticalData(log, conn)

	// getAccountHistory(log, conn)
	// getAccountHistoryWithForEntries(log, conn)
	getSyntheticReportFullPath(log, conn)
}
