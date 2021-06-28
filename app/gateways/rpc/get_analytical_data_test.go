package rpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_GetAnalyticalData(t *testing.T) {
	t.Run("should get analytical data successfully", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetAnalyticalDataFunc: func(ctx context.Context, accountQuery vos.AccountQuery, fn func(vos.Statement) error) error {
				return nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetAnalyticalDataRequest{
			AccountPath: testdata.GenerateAccountPath(),
		}

		grpcServerMocked := &mocks.LedgerService_GetAnalyticalDataServerMock{
			ContextFunc: func() context.Context {
				return context.Background()
			},
			SendFunc: func(getAnalyticalDataResponse *proto.GetAnalyticalDataResponse) error {
				return nil
			},
		}

		err := api.GetAnalyticalData(request, grpcServerMocked)
		assert.NoError(t, err)
	})

	t.Run("should return an error if account query is invalid", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetAnalyticalDataFunc: func(ctx context.Context, accountQuery vos.AccountQuery, fn func(vos.Statement) error) error {
				return nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetAnalyticalDataRequest{
			AccountPath: "liability.bacen.",
		}

		grpcServerMocked := &mocks.LedgerService_GetAnalyticalDataServerMock{
			ContextFunc: func() context.Context {
				return context.Background()
			},
		}

		err := api.GetAnalyticalData(request, grpcServerMocked)
		respStatus, ok := status.FromError(err)

		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, respStatus.Code())
		assert.Equal(t, "invalid account structure", respStatus.Message())
	})
}
