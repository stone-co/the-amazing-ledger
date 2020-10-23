package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
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
	response, err := c.client.SaveTransaction(ctx, transaction.Message)
	if err != nil {
		return fmt.Errorf("save transaction failed: %w", err)
	}

	if response.Error != entities.NoError.Error() {
		return errors.New(response.Error)
	}

	return nil
}
