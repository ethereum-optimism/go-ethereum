-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.ancient_hashes (
  block_number INTEGER PRIMARY KEY,
  hash BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.ancient_hashes;