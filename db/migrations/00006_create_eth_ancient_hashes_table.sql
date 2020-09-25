-- +goose Up
CREATE TABLE IF NOT EXISTS eth.ancient_hashes (
  block_number BIGINT UNIQUE NOT NULL,
  hash BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.ancient_hashes;