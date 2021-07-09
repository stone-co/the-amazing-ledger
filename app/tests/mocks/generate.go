package mocks

//go:generate moq -pkg mocks -out ledger_repository_mock.go ../../domain Repository
//go:generate moq -pkg mocks -out ledger_usecase_mock.go ../../domain UseCase
//go:generate moq -pkg mocks -out grpc_stream_get_account_history_mock.go ../../../gen/ledger LedgerService_GetAccountHistoryServer
//go:generate moq -pkg mocks -out grpc_stream_get_analytical_data_mock.go ../../../gen/ledger LedgerService_GetAnalyticalDataServer
