package rpc

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetAnalyticalData(request *proto.GetAnalyticalDataRequest, stream proto.LedgerService_GetAnalyticalDataServer) error {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAnalyticalData",
	})

	accountPath, err := vos.NewAccountPath(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account path")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	fn := func(st vos.Statement) error {
		if err = stream.Send(&proto.GetAnalyticalDataResponse{
			AccountId: st.Account,
			Operation: proto.Operation(st.Operation),
			Amount:    int32(st.Amount),
		}); err != nil {
			return err
		}

		return nil
	}

	err = a.UseCase.GetAnalyticalData(stream.Context(), *accountPath, fn)
	if err != nil {
		log.WithError(err).Error("can't get account")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
