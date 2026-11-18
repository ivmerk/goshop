-- +goose Up
CREATE TABLE users (
    uuid                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login                TEXT        NOT NULL UNIQUE,
    email                TEXT        NOT NULL,
    password_hash        TEXT        NOT NULL,
    notification_methods JSONB       NOT NULL DEFAULT '[]'::jsonb,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS users;
