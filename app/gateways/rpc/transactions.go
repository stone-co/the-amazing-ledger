package rpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*emptypb.Empty, error) {
	defer newrelic.FromContext(ctx).StartSegment("CreateTransaction").End()

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
		entryID, entryErr := uuid.Parse(entry.Id)
		if entryErr != nil {
			errMsg := "error parsing entry id"
			log.WithError(err).Error(errMsg)
			return nil, status.Error(codes.InvalidArgument, errMsg)
		}

		domainEntry, domainErr := entities.NewEntry(
			entryID,
			vos.OperationType(proto.Operation_value[entry.Operation.String()]),
			entry.AccountId,
			vos.Version(entry.ExpectedVersion),
			int(entry.Amount),
		)
		if domainErr != nil {
			log.WithError(err).Error("error creating entry")
			return nil, status.Error(codes.InvalidArgument, domainErr.Error())
		}

		domainEntries[i] = domainEntry
	}

	competenceDate := time.Unix(req.CompetenceDate.Seconds, 0).UTC()
	if competenceDate.After(time.Now().UTC()) {
		return nil, status.Error(codes.InvalidArgument, "competence date set to the future")
	}

	tx, err := entities.NewTransaction(tid, req.Event, req.Company, competenceDate, domainEntries...)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if err := a.UseCase.CreateTransaction(ctx, tx); err != nil {
		log.WithError(err).Error("creating transaction")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
