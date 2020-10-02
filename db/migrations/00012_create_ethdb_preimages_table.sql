-- +goose Up
CREATE TABLE IF NOT EXISTS ethdb.preimages (
  preimage_key BYTEA PRIMARY KEY,
  preimage BYTEA NOT NULL
);

-- +goose Down
DROP TABLE ethdb.preimages;