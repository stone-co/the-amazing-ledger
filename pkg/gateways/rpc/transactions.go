package rpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (a *API) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*empty.Empty, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "CreateTransaction",
	})

	tid, err := uuid.Parse(req.Id)
	if err != nil {
		errMsg := "error parsing transaction id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}

	domainEntries := make([]entities.Entry, len(req.Entries))
	for i, entry := range req.Entries {
		entryID, err := uuid.Parse(entry.Id)
		if err != nil {
			errMsg := "error parsing entry id"
			log.WithError(err).Error(errMsg)
			return nil, status.Error(codes.InvalidArgument, errMsg)
		}

		domainEntry, err := entities.NewEntry(
			entryID,
			entities.OperationType(proto.Operation_value[entry.Operation.String()]),
			entry.AccountId,
			entities.Version(entry.ExpectedVersion),
			int(entry.Amount),
		)
		if err != nil {
			log.WithError(err).Error("error creating entry")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		domainEntries[i] = *domainEntry
	}

	if err := a.UseCase.CreateTransaction(ctx, tid, domainEntries); err != nil {
		log.WithError(err).Error("creating transaction")
		if err == entities.ErrInvalidVersion {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
