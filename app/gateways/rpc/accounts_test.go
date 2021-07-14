package rpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_GetAccountBalance_Success(t *testing.T) {
	t.Run("should get account balance successfully", func(t *testing.T) {
		accountBalance := vos.AccountBalance{
			CurrentVersion: vos.Version(1),
			TotalCredit:    200,
			TotalDebit:     100,
		}
		mockedUsecase := &mocks.UseCaseMock{
			GetAccountBalanceFunc: func(ctx context.Context, accountPath vos.AccountPath) (vos.AccountBalance, error) {
				accountBalance.Account = accountPath

				return accountBalance, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetAccountBalanceRequest{
			AccountPath: testdata.GenerateAccountPath(),
		}

		got, err := api.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, int64(accountBalance.TotalDebit), got.TotalDebit)
		assert.Equal(t, int64(accountBalance.TotalCredit), got.TotalCredit)
		assert.Equal(t, int64(accountBalance.Balance()), got.Balance)
		assert.Equal(t, accountBalance.Account.Name(), got.AccountPath)
	})
}

func TestAPI_GetAccountBalance_InvalidRequest(t *testing.T) {
	testCases := []struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.GetAccountBalanceRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should return an error if account name is invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.GetAccountBalanceRequest{
				AccountPath: "liability.clients",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: app.ErrInvalidAccountStructure.Error(),
		},
		{
			name: "should return an error if account does not exist",
			useCaseSetup: &mocks.UseCaseMock{
				GetAccountBalanceFunc: func(ctx context.Context, accountPath vos.AccountPath) (vos.AccountBalance, error) {
					return vos.AccountBalance{}, app.ErrAccountNotFound
				},
			},
			request: &proto.GetAccountBalanceRequest{
				AccountPath: testdata.GenerateAccountPath(),
			},
			expectedCode:    codes.NotFound,
			expectedMessage: "account not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			_, err := api.GetAccountBalance(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}

func TestAPI_QueryAggregatedBalance_Success(t *testing.T) {
	t.Run("should get aggregated balance successfully", func(t *testing.T) {
		accountQuery, err := vos.NewAccountQuery("liability.stone.clients.*")
		assert.NoError(t, err)

		queryBalance := vos.NewQueryBalance(accountQuery, 100)
		mockedUsecase := &mocks.UseCaseMock{
			QueryAggregatedBalanceFunc: func(ctx context.Context, accountQuery vos.AccountQuery) (vos.QueryBalance, error) {
				return queryBalance, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.QueryAggregatedBalanceRequest{
			Query: "liability.stone.clients.*",
		}

		got, err := api.QueryAggregatedBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, int64(queryBalance.Balance), got.Balance)
		assert.Equal(t, queryBalance.Query.Value(), got.Query)

	})
}

func TestAPI_QueryAggregatedBalance_InvalidRequest(t *testing.T) {
	testCases := []struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.QueryAggregatedBalanceRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should return an error if account query is invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.QueryAggregatedBalanceRequest{
				Query: "liability.clients.",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: app.ErrInvalidAccountComponentSize.Error(),
		},
		{
			name: "should return an error if account does not exist",
			useCaseSetup: &mocks.UseCaseMock{
				QueryAggregatedBalanceFunc: func(ctx context.Context, accountQuery vos.AccountQuery) (vos.QueryBalance, error) {
					return vos.QueryBalance{}, app.ErrAccountNotFound
				},
			},
			request: &proto.QueryAggregatedBalanceRequest{
				Query: "liability.clients.available",
			},
			expectedCode:    codes.NotFound,
			expectedMessage: "account not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			_, err := api.QueryAggregatedBalance(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}
