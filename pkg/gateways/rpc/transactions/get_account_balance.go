package transactions

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAccountBalance(ctx context.Context, in *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountPath := in.AccountPath

	accountName, err := entities.NewAccountName(accountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := h.UseCase.GetAccountBalance(ctx, *accountName)

	if err != nil {
		if err == entities.ErrAccountNotFound {
			log.WithError(err).Error("account name does not exist")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error("can't get account")
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.GetAccountBalanceResponse{
		AccountPath:    accountBalance.AccountName.Name(),
		CurrentVersion: accountBalance.CurrentVersion.ToUInt64(),
		TotalCredit:    int32(accountBalance.TotalCredit),
		TotalDebit:     int32(accountBalance.TotalDebit),
		Balance:        int32(accountBalance.Balance()),
	}

	return response, nil
}
