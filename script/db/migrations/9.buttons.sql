-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE buttons(
    id TEXT PRIMARY KEY,
    caption TEXT NOT NULL,
    operation TEXT NOT NULL,
    is_delete_after_process BOOLEAN NOT NULL,
    style TEXT NOT NULL,
    url TEXT NOT NULL,
    payload BYTEA,
    pay BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE buttons;
