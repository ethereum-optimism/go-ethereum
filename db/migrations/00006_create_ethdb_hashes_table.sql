-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.hashes (
  hash_key BYTEA PRIMARY KEY,
  hash BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES ethdb.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX hashes_header_fk ON ethdb.hashes USING btree (header_fk);

-- +goose Down
DROP TABLE ethdb.header_hashes;
DROP INDEX ethdb.hashes_header_fk;