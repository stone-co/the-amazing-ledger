package rpc

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"
)

func (a *API) GetSyntheticReport(ctx context.Context, request *proto.GetSyntheticReportRequest) (*proto.GetSyntheticReportResponse, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetSyntheticReport",
	})

	query, err := vos.NewAccountQuery(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("Invalid account query")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	level := int(request.Level) // that's ok to convert int32 to int, since int can be int32 or int64 depending on the used system

	startTime := time.Unix(0, request.StartTime)
	endTime := time.Unix(0, request.EndTime)

	syntheticReport, err := a.UseCase.GetSyntheticReport(ctx, query, level, startTime, endTime)
	if err != nil {
		log.WithError(err).Error("can't get report")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.GetSyntheticReportResponse{
		TotalCredit: int64(syntheticReport.TotalCredit),
		TotalDebit:  int64(syntheticReport.TotalDebit),
		Paths:       toProto(syntheticReport.Paths),
	}, nil
}

func toProto(paths []vos.Path) []*proto.Path {
	protoPaths := []*proto.Path{}

	for _, element := range paths {
		protoPaths = append(protoPaths, &proto.Path{
			Account: element.Account.Name(),
			Credit:  int64(element.Credit),
			Debit:   int64(element.Debit),
		})
	}

	return protoPaths
}
