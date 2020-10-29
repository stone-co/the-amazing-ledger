package ledger

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

type Transaction struct {
	Message *proto.SaveTransactionRequest
}

func (c *Connection) NewTransaction(id uuid.UUID) *Transaction {
	transaction := &Transaction{}
	transaction.Message = &proto.SaveTransactionRequest{
		Id: id.String(),
	}

	return transaction
}

func (c *Connection) SaveTransaction(ctx context.Context, transaction *Transaction) error {
	_, err := c.client.SaveTransaction(ctx, transaction.Message)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return fmt.Errorf(e.Message())
		}

		return fmt.Errorf("not able to parse error returned %v", err)
	}

	return nil
}
