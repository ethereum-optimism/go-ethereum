-- +goose Up
CREATE TABLE IF NOT EXISTS eth.bloom_bits (
  bb_key BYTEA UNIQUE NOT NULL,
  bits BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.bloom_bits;