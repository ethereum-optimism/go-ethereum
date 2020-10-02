-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.bloom_bits (
  bb_key BYTEA PRIMARY KEY,
  bits BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.bloom_bits;