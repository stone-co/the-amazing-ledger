package transactions

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

func (h *Handler) SaveTransaction(ctx context.Context, in *proto.SaveTransactionRequest) (*empty.Empty, error) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "SaveTransaction",
	})

	domainEntries := []entities.Entry{}
	for _, entry := range in.Entries {
		entryID, err := uuid.Parse(entry.Id)
		if err != nil {
			log.WithError(err).Error("parsing entry id")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		var op entities.OperationType
		if entry.Operation == proto.Operation_OPERATION_DEBIT {
			op = entities.DebitOperation
		} else if entry.Operation == proto.Operation_OPERATION_CREDIT {
			op = entities.CreditOperation
		} else {
			op = entities.InvalidOperation
		}

		domainEntry, err := entities.NewEntry(
			entryID,
			op,
			entry.AccountId,
			entities.Version(entry.ExpectedVersion),
			int(entry.Amount),
		)
		if err != nil {
			log.WithError(err).Error("creating entry")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		domainEntries = append(domainEntries, *domainEntry)
	}

	tid, err := uuid.Parse(in.Id)
	if err != nil {
		log.WithError(err).Error("parsing transaction id")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.UseCase.CreateTransaction(ctx, tid, domainEntries); err != nil {
		log.WithError(err).Error("creating transaction")
		if err == entities.ErrInvalidVersion {
			return nil, status.Error(codes.Aborted, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}
