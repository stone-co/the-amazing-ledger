CREATE TYPE account_type AS ENUM (
  'asset',
  'liability',
  'income',
  'expense',
  'equity'
);

CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY NOT NULL,
  "type" account_type NOT NULL,
  owner_id text NOT NULL,
  owner text NOT NULL,
  name text NOT NULL,
  metadata _text NOT NULL,
  balance int NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT positive_balance CHECK (balance >= 0)
);
