package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	defer newrelic.FromContext(ctx).StartSegment("GetAccountBalance").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountName, err := vos.NewAccountName(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := a.UseCase.GetAccountBalance(ctx, *accountName)
	if err != nil {
		if err == app.ErrAccountNotFound {
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

func (a *API) GetAccountHistory(ctx context.Context, request *proto.GetAccountHistoryRequest) (*proto.GetAccountHistoryResponse, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountHistory",
	})

	accountName, err := vos.NewAccountName(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountHistory, err := a.UseCase.GetAccountHistory(ctx, *accountName)
	if err != nil {
		if err == app.ErrAccountNotFound {
			log.WithError(err).Error("account name does not exist")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error("can't get account")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	entriesHistory := make([]*proto.EntryHistory, len(accountHistory.EntriesHistory))

	for i, entryHistory := range accountHistory.EntriesHistory {
		var e proto.EntryHistory
		entriesHistory[i] = &e

		e.Amount = int64(entryHistory.Amount)

		if entryHistory.Operation == vos.CreditOperation {
			e.Operation = proto.Operation_OPERATION_CREDIT
		} else {
			e.Operation = proto.Operation_OPERATION_DEBIT
		}

		timestamp, err := ptypes.TimestampProto(entryHistory.CreatedAt)
		if err != nil {
			log.WithError(err).Error("can't convert time.Time to proto timestamp")
			return nil, err
		}
		e.CreatedAt = timestamp
	}

	return &proto.GetAccountHistoryResponse{
		AccountPath:    accountHistory.Account.Name(),
		TotalCredit:    int64(accountHistory.TotalCredit),
		TotalDebit:     int64(accountHistory.TotalDebit),
		EntriesHistory: entriesHistory,
	}, nil
}
