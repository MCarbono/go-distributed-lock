BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    invoice_id UUID NOT NULL,
    status VARCHAR NOT NULL,
    item_id UUID NOT NULL,
    quantity numeric NOT NULL,
    value numeric NOT NULL,
    total numeric NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMIT;