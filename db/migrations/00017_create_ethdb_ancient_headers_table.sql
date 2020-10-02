-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.ancient_headers (
  block_number INTEGER PRIMARY KEY,
  header BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.ancient_headers;