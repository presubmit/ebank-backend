BEGIN;

CREATE TABLE companies (
    id uuid PRIMARY KEY,
    created_by uuid REFERENCES users (id),
    name text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    is_active boolean NOT NULL DEFAULT TRUE 
);

CREATE TABLE employees (
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users (id),
    company_id uuid NOT NULL REFERENCES companies (id),
    email text,
    role text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    invitation_sent boolean NOT NULL DEFAULT FALSE,
    is_active boolean NOT NULL DEFAULT TRUE
);

COMMIT;