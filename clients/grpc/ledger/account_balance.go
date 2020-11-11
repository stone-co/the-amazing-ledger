package ledger

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/proto"
)

type AccountBalance struct {
	AccountPath    string
	CurrentVersion entities.Version
	Balance        int
}

type AccountRequest struct {
	Message *proto.GetAccountInfoRequest
}

func (c *Connection) NewAccountRequest(id string) *AccountRequest {
	accountRequest := &AccountRequest{}
	accountRequest.Message = &proto.GetAccountInfoRequest{
		AccountId: id,
	}

	return accountRequest
}

func (c *Connection) GetAccountBalance(ctx context.Context, accountRequest *AccountRequest) (*AccountBalance, error) {
	accountBalanceProto, err := c.client.GetAccountBalance(ctx, accountRequest.Message)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("not able to parse error returned %v", err)
	}

	accountBalance := &AccountBalance{
		AccountPath:    accountBalanceProto.AccountId,
		CurrentVersion: entities.Version(accountBalanceProto.CurrentVersion),
		Balance:        int(accountBalanceProto.Balance),
	}

	return accountBalance, nil
}
