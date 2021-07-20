package rpc

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) ListAccountEntries(ctx context.Context, request *proto.ListAccountEntriesRequest) (*proto.ListAccountEntriesResponse, error) {
	defer newrelic.FromContext(ctx).StartSegment("ListAccountEntries").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "ListAccountEntries",
	})

	account, err := vos.NewAnalyticalAccount(request.AccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if request.StartDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start_date must have a value")
	} else if !request.StartDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "start_date must be valid")
	}

	if request.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "end_date must have a value")
	} else if !request.EndDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "end_date must be valid")
	}

	page, err := pagination.NewPage(request.GetPage())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	req := vos.AccountEntryRequest{
		Account:   account,
		StartDate: request.StartDate.AsTime(),
		EndDate:   request.EndDate.AsTime(),
		Page:      page,
	}

	entries, err := a.UseCase.ListAccountEntries(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	protoEntries := make([]*proto.AccountEntry, 0, len(entries.Entries))
	for _, entry := range entries.Entries {
		metadata, err := structpb.NewStruct(entry.Metadata)
		if err != nil {
			log.WithError(err).Error("failed to convert map to structpb")
			return nil, status.Error(codes.Internal, "internal server error")
		}

		protoEntries = append(protoEntries, &proto.AccountEntry{
			Id:             entry.ID.String(),
			Version:        entry.Version.AsInt64(),
			Operation:      proto.Operation(entry.Operation),
			Amount:         int64(entry.Amount),
			Event:          int32(entry.Event),
			CompetenceDate: timestamppb.New(entry.CompetenceDate),
			Metadata:       metadata,
		})
	}

	return &proto.ListAccountEntriesResponse{
		Entries:       protoEntries,
		NextPageToken: entries.NextPage.Tokenize(),
	}, nil
}
