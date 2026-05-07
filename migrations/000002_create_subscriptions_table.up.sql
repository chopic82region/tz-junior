CREATE TABLE IF NOT EXISTS subscriptions (
    id            BIGSERIAL PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name  TEXT NOT NULL,
    price         TEXT NOT NULL,
    payment_time  TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

