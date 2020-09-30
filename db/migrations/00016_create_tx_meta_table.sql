-- +goose Up
CREATE TABLE IF NOT EXISTS eth.tx_meta (
  meta_key BYTEA UNIQUE NOT NULL,
  meta BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.tx_meta;