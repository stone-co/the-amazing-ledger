package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"

	"github.com/golang/protobuf/ptypes"
)

type Transaction struct {
	Message *pb.CreateTransactionRequest
}

func (c *Connection) NewTransaction(id string) *Transaction {

	transaction := &Transaction{}
	transaction.Message = &proto.CreateTransactionRequest{}
	transaction.Message.Id = id

	return transaction
}

func (c *Connection) CreateTransaction(transaction *Transaction) error {
	response, err := c.client.CreateTransaction(context.Background(), transaction.Message)

	if err != nil {
		return fmt.Errorf("create transaction failed: %w", err)
	}

	if response.Status != pb.CreateTransactionResponse_SUCCESS {
		return fmt.Errorf("create transaction failed: %w", response.Status.String())
	}

	return nil
}
