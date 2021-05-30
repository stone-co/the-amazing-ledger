package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	defer newrelic.FromContext(ctx).StartSegment("GetAccountBalance").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountName, err := vos.NewAccountPath(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := a.UseCase.GetAccountBalance(ctx, accountName)
	if err != nil {
		if err == app.ErrAccountNotFound {
			log.WithError(err).Error("account name does not exist")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error("can't get account")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.GetAccountBalanceResponse{
		AccountPath:    accountBalance.Account.Name(),
		CurrentVersion: accountBalance.CurrentVersion.AsInt(),
		TotalCredit:    int64(accountBalance.TotalCredit),
		TotalDebit:     int64(accountBalance.TotalDebit),
		Balance:        int64(accountBalance.Balance()),
	}, nil
}

func (a *API) GetAccountHistory(request *proto.GetAccountHistoryRequest, stream proto.LedgerService_GetAccountHistoryServer) error {
	defer newrelic.FromContext(stream.Context()).StartSegment("GetAccountHistory").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountHistory",
	})

	accountName, err := vos.NewAccountPath(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	fn := func(et vos.EntryHistory) error {
		var timestamp *timestamppb.Timestamp
		timestamp, err = ptypes.TimestampProto(et.CreatedAt)
		if err != nil {
			log.WithError(err).Error("can't convert time.Time to proto timestamp")
			return err
		}

		if err = stream.Send(&proto.GetAccountHistoryResponse{
			Amount:    int64(et.Amount),
			Operation: proto.Operation(et.Operation),
			CreatedAt: timestamp,
		}); err != nil {
			return err
		}

		return nil
	}

	err = a.UseCase.GetAccountHistory(stream.Context(), accountName, fn)
	if err != nil {
		log.WithError(err).Error("can't get account")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
