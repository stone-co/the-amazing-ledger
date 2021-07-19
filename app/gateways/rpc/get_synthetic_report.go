package rpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetSyntheticReport(ctx context.Context, request *proto.GetSyntheticReportRequest) (*proto.GetSyntheticReportResponse, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetSyntheticReport",
	})

	query, err := vos.NewAccount(request.Filters.AccountQuery)
	if err != nil {
		log.WithError(err).Error("Invalid account query")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	level := int(request.Filters.Level) // that's ok to convert int32 to int, since int can be int32 or int64 depending on the used system

	startTime := time.Unix(0, request.Filters.StartTime)
	endTime := time.Unix(0, request.Filters.EndTime)

	syntheticReport, err := a.UseCase.GetSyntheticReport(ctx, query, level, startTime, endTime)
	if err != nil {
		log.WithError(err).Error("can't get synthetic report")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetSyntheticReportResponse{
		TotalCredit: syntheticReport.TotalCredit,
		TotalDebit:  syntheticReport.TotalDebit,
		Paths:       toProto(syntheticReport.Paths),
	}, nil
}

func toProto(paths []vos.Path) []*proto.Path {
	protoPaths := []*proto.Path{}

	for _, element := range paths {
		protoPaths = append(protoPaths, &proto.Path{
			Account: element.Account.Value(),
			Credit:  element.Credit,
			Debit:   element.Debit,
		})
	}

	return protoPaths
}
