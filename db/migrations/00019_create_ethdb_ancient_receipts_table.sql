-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.ancient_receipts (
  block_number INTEGER PRIMARY KEY,
  receipts BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.ancient_receipts;