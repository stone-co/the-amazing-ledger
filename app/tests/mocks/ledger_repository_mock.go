// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"sync"
)

// Ensure, that RepositoryMock does implement domain.Repository.
// If this is not the case, regenerate this file with moq.
var _ domain.Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of domain.Repository.
//
// 	func TestSomethingThatUsesRepository(t *testing.T) {
//
// 		// make and configure a mocked domain.Repository
// 		mockedRepository := &RepositoryMock{
// 			CreateTransactionFunc: func(contextMoqParam context.Context, transaction entities.Transaction) error {
// 				panic("mock out the CreateTransaction method")
// 			},
// 			GetAccountBalanceFunc: func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
// 				panic("mock out the GetAccountBalance method")
// 			},
// 			GetAccountHistoryFunc: func(ctxt context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
// 				panic("mock out the GetAccountHistory method")
// 			},
// 			GetAnalyticalDataFunc: func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
// 				panic("mock out the GetAnalyticalData method")
// 			},
// 			QueryAggregatedBalanceFunc: func(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error) {
// 				panic("mock out the QueryAggregatedBalance method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires domain.Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// CreateTransactionFunc mocks the CreateTransaction method.
	CreateTransactionFunc func(contextMoqParam context.Context, transaction entities.Transaction) error

	// GetAccountBalanceFunc mocks the GetAccountBalance method.
	GetAccountBalanceFunc func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)

	// GetAccountHistoryFunc mocks the GetAccountHistory method.
	GetAccountHistoryFunc func(ctxt context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error

	// GetAnalyticalDataFunc mocks the GetAnalyticalData method.
	GetAnalyticalDataFunc func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error

	// QueryAggregatedBalanceFunc mocks the QueryAggregatedBalance method.
	QueryAggregatedBalanceFunc func(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateTransaction holds details about calls to the CreateTransaction method.
		CreateTransaction []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Transaction is the transaction argument value.
			Transaction entities.Transaction
		}
		// GetAccountBalance holds details about calls to the GetAccountBalance method.
		GetAccountBalance []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Account is the account argument value.
			Account vos.AccountPath
		}
		// GetAccountHistory holds details about calls to the GetAccountHistory method.
		GetAccountHistory []struct {
			// Ctxt is the ctxt argument value.
			Ctxt context.Context
			// Account is the account argument value.
			Account vos.AccountPath
			// Fn is the fn argument value.
			Fn func(vos.EntryHistory) error
		}
		// GetAnalyticalData holds details about calls to the GetAnalyticalData method.
		GetAnalyticalData []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Query is the query argument value.
			Query vos.AccountQuery
			// Fn is the fn argument value.
			Fn func(vos.Statement) error
		}
		// QueryAggregatedBalance holds details about calls to the QueryAggregatedBalance method.
		QueryAggregatedBalance []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Account is the account argument value.
			Account vos.AccountQuery
		}
	}
	lockCreateTransaction      sync.RWMutex
	lockGetAccountBalance      sync.RWMutex
	lockGetAccountHistory      sync.RWMutex
	lockGetAnalyticalData      sync.RWMutex
	lockQueryAggregatedBalance sync.RWMutex
}

// CreateTransaction calls CreateTransactionFunc.
func (mock *RepositoryMock) CreateTransaction(contextMoqParam context.Context, transaction entities.Transaction) error {
	if mock.CreateTransactionFunc == nil {
		panic("RepositoryMock.CreateTransactionFunc: method is nil but Repository.CreateTransaction was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Transaction     entities.Transaction
	}{
		ContextMoqParam: contextMoqParam,
		Transaction:     transaction,
	}
	mock.lockCreateTransaction.Lock()
	mock.calls.CreateTransaction = append(mock.calls.CreateTransaction, callInfo)
	mock.lockCreateTransaction.Unlock()
	return mock.CreateTransactionFunc(contextMoqParam, transaction)
}

// CreateTransactionCalls gets all the calls that were made to CreateTransaction.
// Check the length with:
//     len(mockedRepository.CreateTransactionCalls())
func (mock *RepositoryMock) CreateTransactionCalls() []struct {
	ContextMoqParam context.Context
	Transaction     entities.Transaction
} {
	var calls []struct {
		ContextMoqParam context.Context
		Transaction     entities.Transaction
	}
	mock.lockCreateTransaction.RLock()
	calls = mock.calls.CreateTransaction
	mock.lockCreateTransaction.RUnlock()
	return calls
}

// GetAccountBalance calls GetAccountBalanceFunc.
func (mock *RepositoryMock) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	if mock.GetAccountBalanceFunc == nil {
		panic("RepositoryMock.GetAccountBalanceFunc: method is nil but Repository.GetAccountBalance was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Account vos.AccountPath
	}{
		Ctx:     ctx,
		Account: account,
	}
	mock.lockGetAccountBalance.Lock()
	mock.calls.GetAccountBalance = append(mock.calls.GetAccountBalance, callInfo)
	mock.lockGetAccountBalance.Unlock()
	return mock.GetAccountBalanceFunc(ctx, account)
}

// GetAccountBalanceCalls gets all the calls that were made to GetAccountBalance.
// Check the length with:
//     len(mockedRepository.GetAccountBalanceCalls())
func (mock *RepositoryMock) GetAccountBalanceCalls() []struct {
	Ctx     context.Context
	Account vos.AccountPath
} {
	var calls []struct {
		Ctx     context.Context
		Account vos.AccountPath
	}
	mock.lockGetAccountBalance.RLock()
	calls = mock.calls.GetAccountBalance
	mock.lockGetAccountBalance.RUnlock()
	return calls
}

// GetAccountHistory calls GetAccountHistoryFunc.
func (mock *RepositoryMock) GetAccountHistory(ctxt context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
	if mock.GetAccountHistoryFunc == nil {
		panic("RepositoryMock.GetAccountHistoryFunc: method is nil but Repository.GetAccountHistory was just called")
	}
	callInfo := struct {
		Ctxt    context.Context
		Account vos.AccountPath
		Fn      func(vos.EntryHistory) error
	}{
		Ctxt:    ctxt,
		Account: account,
		Fn:      fn,
	}
	mock.lockGetAccountHistory.Lock()
	mock.calls.GetAccountHistory = append(mock.calls.GetAccountHistory, callInfo)
	mock.lockGetAccountHistory.Unlock()
	return mock.GetAccountHistoryFunc(ctxt, account, fn)
}

// GetAccountHistoryCalls gets all the calls that were made to GetAccountHistory.
// Check the length with:
//     len(mockedRepository.GetAccountHistoryCalls())
func (mock *RepositoryMock) GetAccountHistoryCalls() []struct {
	Ctxt    context.Context
	Account vos.AccountPath
	Fn      func(vos.EntryHistory) error
} {
	var calls []struct {
		Ctxt    context.Context
		Account vos.AccountPath
		Fn      func(vos.EntryHistory) error
	}
	mock.lockGetAccountHistory.RLock()
	calls = mock.calls.GetAccountHistory
	mock.lockGetAccountHistory.RUnlock()
	return calls
}

// GetAnalyticalData calls GetAnalyticalDataFunc.
func (mock *RepositoryMock) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	if mock.GetAnalyticalDataFunc == nil {
		panic("RepositoryMock.GetAnalyticalDataFunc: method is nil but Repository.GetAnalyticalData was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Query vos.AccountQuery
		Fn    func(vos.Statement) error
	}{
		Ctx:   ctx,
		Query: query,
		Fn:    fn,
	}
	mock.lockGetAnalyticalData.Lock()
	mock.calls.GetAnalyticalData = append(mock.calls.GetAnalyticalData, callInfo)
	mock.lockGetAnalyticalData.Unlock()
	return mock.GetAnalyticalDataFunc(ctx, query, fn)
}

// GetAnalyticalDataCalls gets all the calls that were made to GetAnalyticalData.
// Check the length with:
//     len(mockedRepository.GetAnalyticalDataCalls())
func (mock *RepositoryMock) GetAnalyticalDataCalls() []struct {
	Ctx   context.Context
	Query vos.AccountQuery
	Fn    func(vos.Statement) error
} {
	var calls []struct {
		Ctx   context.Context
		Query vos.AccountQuery
		Fn    func(vos.Statement) error
	}
	mock.lockGetAnalyticalData.RLock()
	calls = mock.calls.GetAnalyticalData
	mock.lockGetAnalyticalData.RUnlock()
	return calls
}

// QueryAggregatedBalance calls QueryAggregatedBalanceFunc.
func (mock *RepositoryMock) QueryAggregatedBalance(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error) {
	if mock.QueryAggregatedBalanceFunc == nil {
		panic("RepositoryMock.QueryAggregatedBalanceFunc: method is nil but Repository.QueryAggregatedBalance was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Account vos.AccountQuery
	}{
		Ctx:     ctx,
		Account: account,
	}
	mock.lockQueryAggregatedBalance.Lock()
	mock.calls.QueryAggregatedBalance = append(mock.calls.QueryAggregatedBalance, callInfo)
	mock.lockQueryAggregatedBalance.Unlock()
	return mock.QueryAggregatedBalanceFunc(ctx, account)
}

// QueryAggregatedBalanceCalls gets all the calls that were made to QueryAggregatedBalance.
// Check the length with:
//     len(mockedRepository.QueryAggregatedBalanceCalls())
func (mock *RepositoryMock) QueryAggregatedBalanceCalls() []struct {
	Ctx     context.Context
	Account vos.AccountQuery
} {
	var calls []struct {
		Ctx     context.Context
		Account vos.AccountQuery
	}
	mock.lockQueryAggregatedBalance.RLock()
	calls = mock.calls.QueryAggregatedBalance
	mock.lockQueryAggregatedBalance.RUnlock()
	return calls
}