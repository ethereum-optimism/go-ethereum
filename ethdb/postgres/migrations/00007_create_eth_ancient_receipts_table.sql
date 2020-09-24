-- +goose Up
CREATE TABLE IF NOT EXISTS eth.ancient_receipts (
  block_number BIGINT UNIQUE NOT NULL,
  receipts BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.ancient_receipts;