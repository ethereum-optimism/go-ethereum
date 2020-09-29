-- +goose Up
CREATE TABLE IF NOT EXISTS eth.preimages (
  preimage_key BYTEA UNIQUE NOT NULL,
  preimage BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.preimages;