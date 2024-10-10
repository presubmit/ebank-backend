BEGIN;

CREATE TABLE refresh_tokens (
  id uuid PRIMARY KEY,
  user_id uuid REFERENCES users (id),
  access_token_id uuid,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  expires_at timestamptz NOT NULL
);

COMMIT;