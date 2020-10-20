package client

import (
	"encoding/json"

	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
)

type Operation int32

const (
	DEBIT  Operation = 0
	CREDIT Operation = 1
)

func (t *Transaction) AddEntry(id string, accountId string, expectedVersion uint64, operation Operation, amount int) error {

	var pbOperation pb.Operation = 0

	if operation == DEBIT {
		pbOperation = pb.Operation_DEBIT
	} else {
		pbOperation = pb.Operation_CREDIT
	}

	s.Message.Entries = append(s.Message.Entries, &pb.Entry{
		Id:              id,
		AccountId:       accountId,
		ExpectedVersion: expectedVersion,
		Operation:       pbOperation,
		Amount:          amount,
	})

	return nil
}
