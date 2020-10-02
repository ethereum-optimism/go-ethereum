-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.ancient_tds (
  block_number INTEGER PRIMARY KEY,
  td BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.ancient_tds;