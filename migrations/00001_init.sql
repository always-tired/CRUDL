-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_user_service_start_unique
ON subscriptions (user_id, service_name, start_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions (start_date);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;
