BEGIN;

CREATE TYPE account_type AS ENUM (
	'asset',
	'liability',
	'income',
	'expense',
	'equity'
);

CREATE TABLE accounts (
	id         UUID PRIMARY KEY,
	type       ACCOUNT_TYPE NOT NULL,
	owner_id   TEXT NOT NULL,
	owner      TEXT NOT NULL,
	name       TEXT NOT NULL,
	metadata   TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE account_class AS ENUM (
	'liability',
	'assets',
	'income',
	'expense',
	'equity'
);

CREATE TYPE operation_type AS ENUM (
	'debit',
	'credit'
);

CREATE TABLE entries (
	id               UUID primary key,
	account_class    ACCOUNT_CLASS not null,
	account_group    TEXT not null,
	account_subgroup TEXT not null,
	account_id       TEXT not null,
	operation        OPERATION_TYPE not null,
	amount           INT not null,
	version          BIGINT not null,
	transaction_id   UUID not null,
	created_at       TIMESTAMPTZ not null default CURRENT_TIMESTAMP
);

COMMIT;
