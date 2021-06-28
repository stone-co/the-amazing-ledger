package ledger

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

type EntryHistory struct {
	amount    int
	operation vos.OperationType
	createAt  time.Time
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

func (c *Connection) GetAccountHistory(ctx context.Context, accountPath string) ([]EntryHistory, error) {
	accountRequest := &proto.GetAccountHistoryRequest{
		AccountPath: accountPath,
	}

	stream, err := c.client.GetAccountHistory(ctx, accountRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("%w: %s", ErrUndefined, err)
	}

	entriesHistory := []EntryHistory{}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			if e, ok := status.FromError(err); ok {
				return nil, fmt.Errorf(e.Message())
			}
			return nil, fmt.Errorf("%w: %s", ErrUndefined, err)
		}

		var operation vos.OperationType
		if result.Operation == proto.Operation_OPERATION_CREDIT {
			operation = vos.CreditOperation
		} else {
			operation = vos.DebitOperation
		}

		ok := result.CreatedAt.IsValid()
		if !ok {
			return nil, fmt.Errorf("invalid timestamp received")
		}

		ts := result.CreatedAt.AsTime()
		entriesHistory = append(entriesHistory, EntryHistory{
			amount:    int(result.Amount),
			operation: operation,
			createAt:  ts,
		})
	}

	return entriesHistory, nil
}
