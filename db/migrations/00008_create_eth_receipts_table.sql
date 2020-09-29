-- +goose Up
CREATE TABLE IF NOT EXISTS eth.receipts (
  receipt_key BYTEA UNIQUE NOT NULL,
  receipts BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

-- +goose Down
DROP TABLE eth.receipts;