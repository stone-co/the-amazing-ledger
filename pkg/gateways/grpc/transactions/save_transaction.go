package transactions

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
)

func (h *Handler) SaveTransaction(ctx context.Context, in *pb.SaveTransactionRequest) (*pb.SaveTransactionResponse, error) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "SaveTransaction",
	})

	entries := []entities.Entry{}
	for _, entry := range in.Entries {
		entryID, err := uuid.Parse(entry.Id)
		if err != nil {
			log.WithError(err).Error("parsing entry id")
			return &pb.SaveTransactionResponse{Error: entities.ErrInvalidData.Error()}, nil
		}
		var op entities.OperationType
		if entry.Operation == pb.Operation_DEBIT {
			op = entities.DebitOperation
		} else {
			op = entities.CreditOperation
		}
		entries = append(entries, entities.Entry{
			ID:        entryID,
			Operation: op,
			AccountID: entry.AccountId,
			Version:   entities.Version(entry.ExpectedVersion),
			Amount:    int(entry.Amount),
		})
	}

	tid, err := uuid.Parse(in.Id)
	if err != nil {
		log.WithError(err).Error("parsing transaction id")
		return &pb.SaveTransactionResponse{Error: entities.ErrInvalidData.Error()}, nil
	}

	if err := h.UseCase.CreateTransaction(ctx, tid, entries); err != nil {
		log.WithError(err).Error("creating transaction")
		return &pb.SaveTransactionResponse{Error: err.Error()}, nil
	}

	return &pb.SaveTransactionResponse{Error: entities.NoError.Error()}, nil
}
