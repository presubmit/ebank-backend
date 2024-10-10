BEGIN;

CREATE TABLE counterparties (
  id uuid PRIMARY KEY,
  company_id uuid NOT NULL REFERENCES companies(id),
  first_name text,
  last_name text,
  company_name text,
  type text NOT NULL, -- individual | company
  iban text,
  currency text NOT NULL,
  country text NOT NULL,
  created_by uuid NOT NULL REFERENCES employees(id),
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
  id uuid, 
  leg_id uuid PRIMARY KEY,
  company_id uuid NOT NULL REFERENCES companies (id),
  account_id uuid NOT NULL REFERENCES accounts (id),
  other_account_id uuid REFERENCES accounts (id),
  card_id uuid, -- REFERENCES cards(id),
  counterparty_id uuid REFERENCES counterparties(id),
  amount bigint NOT NULL,
  currency text NOT NULL,
  description text,
  type text NOT NULL, -- transfer, card, refund, exchange
  status text,
  after_balance bigint NOT NULL,
  created_by uuid NOT NULL REFERENCES employees(id),
  created_at timestamptz NOT NULL DEFAULT NOW()
);

COMMIT;