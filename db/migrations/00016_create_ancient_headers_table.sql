-- +goose Up
CREATE TABLE IF NOT EXISTS eth.ancient_headers (
  block_number BIGINT UNIQUE NOT NULL,
  header BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.ancient_headers;