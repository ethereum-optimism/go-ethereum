-- +goose Up
CREATE TABLE IF NOT EXISTS eth.hashes (
  hash_key BYTEA UNIQUE NOT NULL,
  hash BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

-- +goose Down
DROP TABLE eth.header_hashes;