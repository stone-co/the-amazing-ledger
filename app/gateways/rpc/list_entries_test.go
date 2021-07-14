package rpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_ListAccountEntries_Success(t *testing.T) {
	tests := []struct {
		name         string
		useCaseSetup *mocks.UseCaseMock
		request      *proto.ListAccountEntriesRequest
		want         func() (*proto.ListAccountEntriesResponse, error)
	}{
		{
			name: "should succeed when listing account entries - no cursor",
			useCaseSetup: &mocks.UseCaseMock{
				ListAccountEntriesFunc: func(_ context.Context, _ vos.AccountEntryRequest) (vos.AccountEntryResponse, error) {
					return vos.AccountEntryResponse{
						Entries:  []vos.AccountEntry{},
						NextPage: nil,
					}, nil
				},
			},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
				Page:        nil,
			},
			want: func() (*proto.ListAccountEntriesResponse, error) {
				return &proto.ListAccountEntriesResponse{
					Entries:       []*proto.AccountEntry{},
					NextPageToken: "",
				}, nil
			},
		},
		{
			name: "should succeed when listing account entries - with cursor",
			useCaseSetup: &mocks.UseCaseMock{
				ListAccountEntriesFunc: func(_ context.Context, _ vos.AccountEntryRequest) (vos.AccountEntryResponse, error) {
					cursor, err := pagination.NewCursor(map[string]interface{}{"abc": 123})
					if err != nil {
						return vos.AccountEntryResponse{}, err
					}

					return vos.AccountEntryResponse{
						Entries:  []vos.AccountEntry{},
						NextPage: cursor,
					}, nil
				},
			},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
				Page:        nil,
			},
			want: func() (*proto.ListAccountEntriesResponse, error) {
				cursor, err := pagination.NewCursor(map[string]interface{}{"abc": 123})
				if err != nil {
					return nil, err
				}

				return &proto.ListAccountEntriesResponse{
					Entries:       []*proto.AccountEntry{},
					NextPageToken: cursor.Tokenize(),
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			got, err := api.ListAccountEntries(context.Background(), tt.request)
			assert.NoError(t, err)

			want, err := tt.want()
			assert.NoError(t, err)

			assert.Equal(t, want, got)
			assert.Len(t, tt.useCaseSetup.ListAccountEntriesCalls(), 1)

			account, _ := vos.NewAccountPath(tt.request.AccountPath)
			page, _ := pagination.NewPage(nil)
			assert.Equal(t, vos.AccountEntryRequest{
				Account:   account,
				StartDate: tt.request.StartDate.AsTime(),
				EndDate:   tt.request.EndDate.AsTime(),
				Page:      page,
			}, tt.useCaseSetup.ListAccountEntriesCalls()[0].AccountEntryRequest)
		})
	}
}

func TestAPI_ListAccountEntries_InvalidRequest(t *testing.T) {
	tests := []struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.ListAccountEntriesRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should return an error with an invalid account path",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.$.account1",
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
				Page:        nil,
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "only alphanumeric and underscore characters are supported",
		},
		{
			name:         "should return an error with nil start date",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   nil,
				EndDate:     timestamppb.Now(),
				Page:        nil,
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "start_date must have a value",
		},
		{
			name:         "should return an error with invalid start date",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate: &timestamppb.Timestamp{
					Seconds: -100,
					Nanos:   -100,
				},
				EndDate: timestamppb.Now(),
				Page:    nil,
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "start_date must be valid",
		},
		{
			name:         "should return an error with nil end date",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   timestamppb.Now(),
				EndDate:     nil,
				Page:        nil,
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "end_date must have a value",
		},
		{
			name:         "should return an error with invalid end date",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   timestamppb.Now(),
				EndDate: &timestamppb.Timestamp{
					Seconds: -100,
					Nanos:   -100,
				},
				Page: nil,
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "end_date must be valid",
		},
		{
			name:         "should return an error with invalid pagination",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability.credit_card.account1",
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
				Page: &proto.RequestPagination{
					PageSize:  0,
					PageToken: "",
				},
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: "invalid page size",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			_, err := api.ListAccountEntries(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}
