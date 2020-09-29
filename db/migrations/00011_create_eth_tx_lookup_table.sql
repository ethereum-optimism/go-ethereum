-- +goose Up
CREATE TABLE IF NOT EXISTS eth.tx_lookups (
  lookup_key BYTEA UNIQUE NOT NULL,
  lookup BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.tx_lookups;