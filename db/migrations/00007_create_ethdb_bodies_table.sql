-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.bodies (
  body_key BYTEA PRIMARY KEY,
  body BYTEA NOT NULL,
  header_fk BYTEA NOT NULL REFERENCES ethdb.headers (header_key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX bodies_header_fk ON ethdb.bodies USING btree (header_fk);

-- +goose Down
DROP TABLE ethdb.block_bodies;
DROP INDEX ethdb.bodies_header_fk;