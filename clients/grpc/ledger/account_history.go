package ledger

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/status"
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

		time, err := ptypes.Timestamp(result.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: can't convert time.Time to proto timestamp", err)

		}
		entriesHistory = append(entriesHistory, EntryHistory{
			amount:    int(result.Amount),
			operation: operation,
			createAt:  time,
		})
	}

	return entriesHistory, nil
}
