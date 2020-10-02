-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.tx_meta (
  meta_key BYTEA PRIMARY KEY,
  meta BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.tx_meta;