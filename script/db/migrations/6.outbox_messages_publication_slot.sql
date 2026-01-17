-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- use only with debezium config
-- SELECT pg_create_logical_replication_slot('postgres_debezium', 'pgoutput');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
