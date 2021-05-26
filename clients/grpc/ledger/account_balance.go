package ledger

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

type AccountBalance struct {
	account        vos.AccountPath
	currentVersion vos.Version
	totalCredit    int
	totalDebit     int
	balance        int
}

func (a AccountBalance) Account() vos.AccountPath {
	return a.account
}

func (a AccountBalance) TotalCredit() int {
	return a.totalCredit
}

func (a AccountBalance) CurrentVersion() vos.Version {
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

		return nil, fmt.Errorf("%w: %s", ErrUndefined, err)
	}

	account, err := vos.NewAccountPath(response.AccountPath)
	if err != nil {
		return nil, err
	}

	accountBalance := &AccountBalance{
		account:        account,
		currentVersion: vos.Version(response.CurrentVersion),
		totalCredit:    int(response.TotalCredit),
		totalDebit:     int(response.TotalDebit),
		balance:        int(response.Balance),
	}

	return accountBalance, nil
}
