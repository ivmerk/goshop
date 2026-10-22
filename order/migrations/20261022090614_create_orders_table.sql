-- +goose Up
CREATE TABLE orders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user"      TEXT          NOT NULL,
    parts       TEXT[]        NOT NULL,
    total       NUMERIC(12,2) NOT NULL CHECK (total >= 0),
    payment     TEXT          NOT NULL,
    transaction TEXT,
    status      TEXT          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS orders;
