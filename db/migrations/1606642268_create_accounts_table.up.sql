BEGIN;

CREATE TABLE accounts (
  id uuid PRIMARY KEY,
  company_id uuid REFERENCES companies (id),
  balance bigint DEFAULT 0,
  currency text NOT NULL,
  name text,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

COMMIT;