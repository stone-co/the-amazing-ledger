package transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {
	target := "/transactions"
	body, _ := json.Marshal(CreateTransactionRequest{})

	t.Run("should return 201 when transaction is successfully created", func(t *testing.T) {
		// prepare
		handler := createHandler(nil)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))

		// test
		http.HandlerFunc(handler.Create).ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusCreated, response.Code)
		// assert.Equal(t, ContentTypeJSON, response.Header().Get("content-type"))
	})

	t.Run("should return 400 when transaction is invalid", func(t *testing.T) {
		// prepare
		handler := createHandler(entities.ErrInvalidBalance)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))

		// test
		http.HandlerFunc(handler.Create).ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		// assert.Equal(t, ContentTypeJSON, response.Header().Get("content-type"))
	})
}

func createHandler(err error) Handler {
	return Handler{
		log: logrus.New(),
		UseCase: &ledger.TransactionsMock{
			OnCreateTransaction: func(ctx context.Context, id uuid.UUID, createdAt time.Time, entries []entities.Entry) error {
				return err
			},
		},
	}
}
