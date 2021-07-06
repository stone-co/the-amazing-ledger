package rpc

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_CreateTransaction_Success(t *testing.T) {
	tests := []*struct {
		name         string
		useCaseSetup *mocks.UseCaseMock
		request      *proto.CreateTransactionRequest
	}{
		{
			name: "should succeed when create a transaction",
			useCaseSetup: &mocks.UseCaseMock{
				CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) error {
					return nil
				},
			},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"requestID": {
									Kind: &structpb.Value_StringValue{
										StringValue: "my-request-id-1",
									},
								},
							},
						},
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"requestID": {
									Kind: &structpb.Value_StringValue{
										StringValue: "my-request-id-2",
									},
								},
							},
						},
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			got, err := api.CreateTransaction(context.Background(), tt.request)
			assert.NoError(t, err)
			assert.Equal(t, &emptypb.Empty{}, got)
		})
	}
}

func TestAPI_CreateTransaction_InvalidRequest(t *testing.T) {
	tests := []*struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.CreateTransactionRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should not create transaction when invalid ID",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: "invalid UUID",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "error parsing transaction id",
		},
		{
			name:         "should not create transaction when invalid entry ID",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              "invalid-entry-id",
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "error parsing entry id",
		},
		{
			name:         "should not create transaction when invalid operation",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_UNSPECIFIED,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "invalid operation",
		},
		{
			name:         "should not create transaction when invalid amount",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          -3,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "invalid amount",
		},
		{
			name:         "should not create transaction when number of entries is less than two",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          100,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode:    codes.Aborted,
			expectedMessage: "invalid entries number",
		},
		{
			name:         "should not create transaction when account is invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       "assets",
						ExpectedVersion: 2,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: app.ErrInvalidAccountStructure.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			_, err := api.CreateTransaction(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}
