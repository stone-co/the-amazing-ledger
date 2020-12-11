package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func defineTransactionWithThreeEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting defineTransactionWithThreeEntries")
	defer log.Println("finishing defineTransactionWithThreeEntries")

	// Define a new transaction with 3 entries
	t := conn.NewTransaction(uuid.New())

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()
	accountID3 := "liability:clients:available:" + uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, vos.NewAccountVersion, vos.CreditOperation, 10000)
	t.AddEntry(uuid.New(), accountID3, vos.NewAccountVersion, vos.CreditOperation, 5000)

	// Save transaction
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)
}
