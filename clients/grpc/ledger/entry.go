package ledger

import (
	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

func (t *Transaction) AddEntry(id uuid.UUID, accountId string, expectedVersion entities.Version, operation entities.OperationType, amount int) {
	var pbOperation proto.Operation

	if operation == entities.DebitOperation {
		pbOperation = proto.Operation_OPERATION_DEBIT
	} else {
		pbOperation = proto.Operation_OPERATION_CREDIT
	}

	t.Message.Entries = append(t.Message.Entries, &proto.Entry{
		Id:              id.String(),
		AccountId:       accountId,
		ExpectedVersion: uint64(expectedVersion),
		Operation:       pbOperation,
		Amount:          int32(amount),
	})
}
