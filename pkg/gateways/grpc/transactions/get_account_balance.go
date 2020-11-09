package transactions

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAccountBalance(ctx context.Context, in *proto.GetAccountInfoRequest) (*proto.GetAccountInfoResponse, error) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountID := in.AccountId

	accountInfo, err := h.UseCase.GetAccountInfo(ctx, accountID)

	if err != nil {
		if err == entities.ErrNotFound {
			log.WithError(err).Error("object not found")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error(" error getting account info")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accountBalance := &proto.GetAccountInfoResponse{
		AccountId:      accountInfo.AccountID,
		CurrentVersion: int64(accountInfo.CurrentVersion.Current()),
		Balance:        int32(accountInfo.Balance()),
	}

	return accountBalance, nil
}
