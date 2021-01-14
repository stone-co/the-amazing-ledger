package ledger

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/status"
)

type AccountHistory struct {
	accountName    vos.AccountName
	totalCredit    int
	totalDebit     int
	entriesHistory []EntryHistory
}

type EntryHistory struct {
	amount    int
	operation vos.OperationType
	createAt  time.Time
}

func (a AccountHistory) AccountName() vos.AccountName {
	return a.accountName
}

func (a AccountHistory) TotalCredit() int {
	return a.totalCredit
}

func (a AccountHistory) TotalDebit() int {
	return a.totalDebit
}

func (a AccountHistory) EntriesHistory() []EntryHistory {
	return a.entriesHistory
}

func (e EntryHistory) Amount() int {
	return e.amount
}

func (e EntryHistory) Operation() vos.OperationType {
	return e.operation
}

func (e EntryHistory) CreateAt() time.Time {
	return e.createAt
}

func (c *Connection) GetAccountHistory(ctx context.Context, accountPath string) (*AccountHistory, error) {
	accountRequest := &proto.GetAccountHistoryRequest{
		AccountPath: accountPath,
	}

	response, err := c.client.GetAccountHistory(ctx, accountRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("%w: %s", ErrUndefined, err)
	}

	accountName, err := vos.NewAccountName(response.AccountPath)
	if err != nil {
		return nil, err
	}

	entriesHistory := make([]EntryHistory, len(response.GetEntriesHistory()))

	for i, entryHistory := range response.GetEntriesHistory() {
		entriesHistory[i].amount = int(entryHistory.Amount)

		if entryHistory.GetOperation() == proto.Operation_OPERATION_CREDIT {
			entriesHistory[i].operation = vos.CreditOperation
		} else {
			entriesHistory[i].operation = vos.DebitOperation
		}

		time, err := ptypes.Timestamp(entryHistory.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: can't convert time.Time to proto timestamp", err)

		}
		entriesHistory[i].createAt = time
	}

	accountHistory := &AccountHistory{
		accountName:    *accountName,
		totalCredit:    int(response.TotalCredit),
		totalDebit:     int(response.TotalDebit),
		entriesHistory: entriesHistory,
	}

	return accountHistory, nil
}
