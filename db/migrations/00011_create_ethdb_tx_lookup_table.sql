-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.tx_lookups (
  lookup_key BYTEA PRIMARY KEY,
  lookup BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.tx_lookups;