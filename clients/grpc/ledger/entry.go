package ledger

import (
	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

type Operation int

const (
	Debit  Operation = 0
	Credit Operation = 1
)

func (t *Transaction) AddEntry(id uuid.UUID, accountId string, expectedVersion entities.Version, operation Operation, amount int) {
	var pbOperation proto.Operation

	if operation == Debit {
		pbOperation = proto.Operation_DEBIT
	} else {
		pbOperation = proto.Operation_CREDIT
	}

	t.Message.Entries = append(t.Message.Entries, &proto.Entry{
		Id:              id.String(),
		AccountId:       accountId,
		ExpectedVersion: uint64(expectedVersion),
		Operation:       pbOperation,
		Amount:          int32(amount),
	})
}
