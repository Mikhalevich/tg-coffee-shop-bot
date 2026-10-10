-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE daily_positions(
    day DATE PRIMARY KEY,
    position INTEGER NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE daily_positions;
