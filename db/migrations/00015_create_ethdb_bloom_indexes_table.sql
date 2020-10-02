-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.bloom_indexes (
  bbi_key BYTEA PRIMARY KEY,
  index BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.bloom_indexes;