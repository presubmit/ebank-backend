BEGIN;

CREATE TABLE cards (
    id uuid PRIMARY KEY,
    employee_id uuid REFERENCES employees (id),
    company_id uuid REFERENCES companies (id),
    created_by uuid REFERENCES employees (id),
    brand text NOT NULL,
    number text NOT NULL,
    expiration_month smallint NOT NULL,
    expiration_year smallint NOT NULL,
    security_code text NOT NULL,
    type text NOT NULL, 
    frozen_at timestamptz DEFAULT NULL, 
    closed_at timestamptz DEFAULT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

COMMIT;