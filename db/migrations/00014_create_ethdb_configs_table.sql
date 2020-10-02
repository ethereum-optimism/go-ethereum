-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.configs (
  config_key BYTEA PRIMARY KEY,
  config BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.configs;