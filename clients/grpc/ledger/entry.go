package ledger

import (
	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (t *Transaction) AddEntry(id uuid.UUID, accountId string, expectedVersion vos.Version, operation vos.OperationType, amount int) {
	var pbOperation proto.Operation

	if operation == vos.DebitOperation {
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
