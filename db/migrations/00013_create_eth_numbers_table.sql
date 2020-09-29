-- +goose Up
CREATE TABLE IF NOT EXISTS eth.numbers (
  number_key BYTEA UNIQUE NOT NULL,
  number BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

-- +goose Down
DROP TABLE eth.numbers;