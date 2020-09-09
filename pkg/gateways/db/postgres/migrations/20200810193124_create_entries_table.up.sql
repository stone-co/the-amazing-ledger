create table if not exists entries (
    id uuid primary key,
	account_id text not null,
	transaction_id uuid not null,
	request_id text null,
	amount int not null,
	balance_after int null,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz null
);