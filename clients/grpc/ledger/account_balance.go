package ledger

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

type AccountBalance struct {
	AccountID      string
	CurrentVersion entities.Version
	Balance        int
}

type AccountID struct {
	Message *proto.GetAccountInfoRequest
}

func (c *Connection) NewAccountID(id string) *AccountID {
	accountID := &AccountID{}
	accountID.Message = &proto.GetAccountInfoRequest{
		AccountId: id,
	}

	return accountID
}

func (c *Connection) GetAccountBalance(ctx context.Context, accountID *AccountID) (*AccountBalance, error) {
	accountBalanceProto, err := c.client.GetAccountBalance(ctx, accountID.Message)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("not able to parse error returned %v", err)
	}

	accountBalance := &AccountBalance{
		AccountID:      accountBalanceProto.AccountId,
		CurrentVersion: entities.Version(accountBalanceProto.CurrentVersion),
		Balance:        int(accountBalanceProto.Balance),
	}

	return accountBalance, nil
}
