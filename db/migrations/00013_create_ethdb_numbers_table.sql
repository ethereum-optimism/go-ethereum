-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.numbers (
  number_key BYTEA PRIMARY KEY,
  number BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES ethdb.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX numbers_header_fk ON ethdb.numbers USING btree (header_fk);

-- +goose Down
DROP TABLE ethdb.numbers;
DROP INDEX ethdb.numbers_header_fk;