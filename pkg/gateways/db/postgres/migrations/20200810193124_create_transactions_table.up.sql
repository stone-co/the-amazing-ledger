create type transaction_type as enum (
  'credit',
  'debit'
);


create table if not exists transactions (
    id uuid primary key,
	account_id text not null,
	operation_id text not null,
	request_id text null,
	"type" transaction_type not null,
	amount int not null,
	balance_after int null,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz null
);