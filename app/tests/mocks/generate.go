package mocks

//go:generate moq -pkg mocks -out ledger_instrumentation_mock.go ../../domain Instrumentation
//go:generate moq -pkg mocks -out ledger_repository_mock.go ../../domain Repository
//go:generate moq -pkg mocks -out ledger_usecase_mock.go ../../domain UseCase
