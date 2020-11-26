package ledger

import (
	"context"
	"fmt"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

type Transaction struct {
	Message *proto.CreateTransactionRequest
}

func (c *Connection) NewTransaction(id uuid.UUID) *Transaction {
	transaction := &Transaction{}
	transaction.Message = &proto.CreateTransactionRequest{
		Id: id.String(),
	}

	return transaction
}

func (c *Connection) SaveTransaction(ctx context.Context, transaction *Transaction) error {
	_, err := c.client.CreateTransaction(ctx, transaction.Message)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return fmt.Errorf(e.Message())
		}

		return fmt.Errorf("not able to parse error returned %v", err)
	}

	return nil
}
