CREATE TYPE account_type AS ENUM (
  'asset',
  'liability'
);

CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY NOT NULL,
  owner text NOT NULL,
  name text NOT NULL,
  owner_id text NOT NULL,
  "type" account_type NOT NULL,
  metadata _text NOT NULL,
  balance int NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updates_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT positive_balance CHECK (balance > 0)
);
