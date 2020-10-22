package ledger

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
)

type Operation int32

const (
	Debit  Operation = 0
	Credit Operation = 1
)

func (t *Transaction) AddEntry(id uuid.UUID, accountId string, expectedVersion entities.Version, operation Operation, amount int) error {
	var pbOperation pb.Operation = 0

	if operation == Debit {
		pbOperation = pb.Operation_DEBIT
	} else {
		pbOperation = pb.Operation_CREDIT
	}

	t.Message.Entries = append(t.Message.Entries, &pb.Entry{
		Id:              id.String(),
		AccountId:       accountId,
		ExpectedVersion: uint64(expectedVersion),
		Operation:       pbOperation,
		Amount:          int32(amount),
	})

	return nil
}
