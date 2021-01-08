package ledger

import (
	"context"
	"fmt"
	"io"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/status"
)

type Statement struct {
	Account   string
	Operation vos.OperationType
	Amount    int
}

func (c *Connection) GetAnalyticalData(ctx context.Context, path string) ([]Statement, error) {

	request := &proto.GetAnalyticalDataRequest{
		AccountPath: path,
	}

	stream, err := c.client.GetAnalyticalData(ctx, request)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("%w: %s", ErrUndefined, err)
	}

	statements := []Statement{}
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

		var op vos.OperationType
		if result.Operation == proto.Operation_OPERATION_DEBIT {
			op = vos.DebitOperation
		} else {
			op = vos.CreditOperation
		}

		statements = append(statements, Statement{
			Account:   result.AccountId,
			Operation: op,
			Amount:    int(result.Amount),
		})
	}

	return statements, nil
}
