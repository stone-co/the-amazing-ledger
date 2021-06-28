package rpc

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetAnalyticalData(request *proto.GetAnalyticalDataRequest, stream proto.LedgerService_GetAnalyticalDataServer) error {
	defer newrelic.FromContext(stream.Context()).StartSegment("GetAnalyticalData").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAnalyticalData",
	})

	query, err := vos.NewAccountQuery(request.AccountPath)
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

	err = a.UseCase.GetAnalyticalData(stream.Context(), query, fn)
	if err != nil {
		log.WithError(err).Error("can't get account")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
