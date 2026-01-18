-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE outbox_answer_payment(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    payment_id TEXT NOT NULL,
    ok BOOLEAN NOT NULL,
    error_msg TEXT NOT NULL,
    is_dispatched BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    dispatched_at TIMESTAMPTZ
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE outbox_payment_answer;
