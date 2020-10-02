-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.tds (
  td_key BYTEA PRIMARY KEY,
  td BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES ethdb.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX tds_header_fk ON ethdb.tds USING btree (header_fk);

-- +goose Down
DROP TABLE ethdb.tds;
DROP INDEX ethdb.tds_header_fk;