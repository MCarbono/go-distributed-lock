BEGIN;

CREATE TABLE IF NOT EXISTS invoices (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    status VARCHAR NOT NULL,
    total numeric,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMIT;