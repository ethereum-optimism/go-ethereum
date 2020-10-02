-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.receipts (
  receipt_key BYTEA PRIMARY KEY,
  receipts BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES ethdb.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX receipts_header_fk ON ethdb.receipts USING btree (header_fk);

-- +goose Down
DROP TABLE ethdb.receipts;
DROP INDEX ethdb.receipts_header_fk;