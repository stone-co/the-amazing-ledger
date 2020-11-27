package rpc

import (
	"context"

	"github.com/sirupsen/logrus"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) CreateAccount(_ context.Context, req *proto.CreateAccountRequest) (*proto.Account, error) {

	log := a.log.WithFields(logrus.Fields{
		"handler": "CreateAccount",
	})

	accountInput := ledger.AccountInput{
		Type:     req.Type,
		OwnerID:  req.OwnerId,
		Owner:    req.Owner,
		Name:     req.Name,
		Metadata: req.Metadata,
	}
	acc, err := a.AccountsUseCase.CreateAccount(accountInput)

	if err != nil {
		log.WithError(err).Error("creating transaction")
		if err == entities.ErrInvalidVersion {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.Account{
		Id: acc.ID,
	}, nil

}

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountName, err := entities.NewAccountName(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := a.TransactionUseCase.GetAccountBalance(ctx, *accountName)
	if err != nil {
		if err == entities.ErrAccountNotFound {
			log.WithError(err).Error("account name does not exist")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error("can't get account")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.GetAccountBalanceResponse{
		AccountPath:    accountBalance.AccountName.Name(),
		CurrentVersion: accountBalance.CurrentVersion.ToUInt64(),
		TotalCredit:    int64(accountBalance.TotalCredit),
		TotalDebit:     int64(accountBalance.TotalDebit),
		Balance:        int64(accountBalance.Balance()),
	}, nil
}
