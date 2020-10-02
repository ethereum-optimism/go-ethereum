-- +goose Up
CREATE TABLE IF NOT EXISTS eth.bodies (
  body_key BYTEA UNIQUE NOT NULL,
  body BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES eth.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX bodies_header_fk ON eth.bodies USING btree (header_fk);

-- +goose Down
DROP TABLE eth.block_bodies;
DROP INDEX eth.bodies_header_fk;