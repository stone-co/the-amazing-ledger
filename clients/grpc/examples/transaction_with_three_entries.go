package main

import (
	"log"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func defineTransactionWithThreeEntries(conn *ledger.Connection) {
	log.Println("defineTransactionWithThreeEntries example starting...")
	defer log.Println("defineTransactionWithThreeEntries example finishing...")

	// Define a new transaction with 3 entries
	t := conn.NewTransaction(uuid.New())

	accountID1 := uuid.New().String()
	accountID2 := uuid.New().String()
	accountID3 := uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, ledger.Debit, 15000)
	t.AddEntry(uuid.New(), accountID2, entities.NewAccountVersion, ledger.Credit, 10000)
	t.AddEntry(uuid.New(), accountID3, entities.NewAccountVersion, ledger.Credit, 5000)

	// Save transaction
	err := conn.SaveTransaction(t)
	AssertEqual(nil, err)
}
