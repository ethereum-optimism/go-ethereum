-- +goose Up
CREATE TABLE IF NOT EXISTS eth.tds (
  td_key BYTEA UNIQUE NOT NULL,
  td BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX tds_header_fk ON eth.tds USING btree (header_fk);

-- +goose Down
DROP TABLE eth.tds;
DROP INDEX eth.tds_header_fk;