package ledger

import (
	"context"
	"fmt"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type AccountBalance struct {
	accountName    entities.AccountName
	currentVersion entities.Version
	totalCredit    int
	totalDebit     int
	balance        int
}

func (a AccountBalance) AccountName() entities.AccountName {
	return a.accountName
}

func (a AccountBalance) TotalCredit() int {
	return a.totalCredit
}

func (a AccountBalance) CurrentVersion() entities.Version {
	return a.currentVersion
}

func (a AccountBalance) TotalDebit() int {
	return a.totalDebit
}

func (a AccountBalance) Balance() int {
	return a.balance
}

func (c *Connection) GetAccountBalance(ctx context.Context, accountPath string) (*AccountBalance, error) {

	accountRequest := &proto.GetAccountBalanceRequest{
		AccountPath: accountPath,
	}

	response, err := c.client.GetAccountBalance(ctx, accountRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, ErrUnmapped.cause(err)
	}

	accountName, err := entities.NewAccountName(response.AccountPath)
	if err != nil {
		return nil, err
	}

	accountBalance := &AccountBalance{
		accountName:    *accountName,
		currentVersion: entities.Version(response.CurrentVersion),
		totalCredit:    int(response.TotalCredit),
		totalDebit:     int(response.TotalDebit),
		balance:        int(response.Balance),
	}

	return accountBalance, nil
}
