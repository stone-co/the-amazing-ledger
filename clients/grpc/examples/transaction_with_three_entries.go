package main

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func defineTransactionWithThreeEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting defineTransactionWithThreeEntries")
	defer log.Println("finishing defineTransactionWithThreeEntries")

	// Define a new transaction with 3 entries
	t := conn.NewTransaction(uuid.New())

	accountID1 := uuid.New().String()
	accountID2 := uuid.New().String()
	accountID3 := uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, entities.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, entities.NewAccountVersion, entities.CreditOperation, 10000)
	t.AddEntry(uuid.New(), accountID3, entities.NewAccountVersion, entities.CreditOperation, 5000)

	// Save transaction
	err := conn.SaveTransaction(t)
	AssertEqual(nil, err)
}
