-- +goose Up
CREATE TABLE IF NOT EXISTS eth.bloom_indexes (
  bbi_key BYTEA UNIQUE NOT NULL,
  index BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.bloom_indexes;