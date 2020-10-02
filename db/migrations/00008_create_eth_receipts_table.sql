-- +goose Up
CREATE TABLE IF NOT EXISTS eth.receipts (
  receipt_key BYTEA UNIQUE NOT NULL,
  receipts BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX receipts_header_fk ON eth.receipts USING btree (header_fk);

-- +goose Down
DROP TABLE eth.receipts;
DROP INDEX eth.receipts_header_fk;