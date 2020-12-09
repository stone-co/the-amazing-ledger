package rpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/mocks"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAPI_CreateTransaction_Success(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	tests := []*struct {
		name    string
		request proto.CreateTransactionRequest
	}{
		{
			name: "should succeed when create a transaction",
			request: proto.CreateTransactionRequest{
				Id: ValidTransactionID,
				Entries: []*proto.Entry{
					{
						Id:              "f6162a96-efa3-4d8b-8636-851a9c1a2cd4",
						AccountId:       ValidAccountID,
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				log:     logger,
				UseCase: mocks.SuccessfulTransactionMock(),
			}
			actual, err := a.CreateTransaction(ctx, &tt.request)
			require.NoError(t, err)
			require.Equal(t, actual, &empty.Empty{})
		})
	}
}

func TestAPI_CreateTransaction_InvalidRequest(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	tests := []*struct {
		name            string
		request         proto.CreateTransactionRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name: "should not create transaction when invalid ID",
			request: proto.CreateTransactionRequest{
				Id: "invalid UUID",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "error parsing transaction id",
		},
		{
			name: "should not create transaction when invalid entry ID",
			request: proto.CreateTransactionRequest{
				Id: ValidTransactionID,
				Entries: []*proto.Entry{
					{
						Id:              "invalid-entry-id",
						AccountId:       ValidAccountID,
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
				},
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "error parsing entry id",
		},
		{
			name: "should not create transaction when invalid operation",
			request: proto.CreateTransactionRequest{
				Id: ValidTransactionID,
				Entries: []*proto.Entry{
					{
						Id:              "f6162a96-efa3-4d8b-8636-851a9c1a2cd4",
						AccountId:       ValidAccountID,
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_UNSPECIFIED,
						Amount:          123,
					},
				},
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "invalid data: operation",
		},
		{
			name: "should not create transaction when invalid amount",
			request: proto.CreateTransactionRequest{
				Id: ValidTransactionID,
				Entries: []*proto.Entry{
					{
						Id:              "f6162a96-efa3-4d8b-8636-851a9c1a2cd4",
						AccountId:       ValidAccountID,
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          -3,
					},
				},
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "invalid data: amount",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				log:     logger,
				UseCase: mocks.SuccessfulTransactionMock(),
			}
			_, err := a.CreateTransaction(ctx, &tt.request)
			respStatus, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, tt.expectedCode.String(), respStatus.Code().String())
			require.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}
