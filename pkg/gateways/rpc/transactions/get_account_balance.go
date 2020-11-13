package transactions

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/rpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAccountBalance(ctx context.Context, in *proto.GetAccountInfoRequest) (*proto.GetAccountInfoResponse, error) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountPath := in.AccountId

	accountName, err := entities.NewAccountName(accountPath)

	accountBalance, err := h.UseCase.GetAccountBalance(ctx, *accountName)

	if err != nil {
		if err == entities.ErrNotFound {
			log.WithError(err).Error("object not found")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error(" error getting account info")
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.GetAccountInfoResponse{
		AccountId:      accountBalance.AccountName.Name(),
		CurrentVersion: int64(accountBalance.CurrentVersion.Current()),
		Balance:        int32(accountBalance.Balance()),
	}

	return response, nil
}
