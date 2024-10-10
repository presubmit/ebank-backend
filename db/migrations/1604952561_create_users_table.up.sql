BEGIN;

CREATE TABLE users (
  id uuid PRIMARY KEY,
  first_name text,
  last_name text,
  email text UNIQUE NOT NULL,
  password text,
  phone_number text,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

COMMIT;

