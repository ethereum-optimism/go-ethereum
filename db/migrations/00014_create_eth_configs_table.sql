-- +goose Up
CREATE TABLE IF NOT EXISTS eth.configs (
  config_key BYTEA UNIQUE NOT NULL,
  config BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.configs;