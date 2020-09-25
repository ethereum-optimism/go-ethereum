-- +goose Up
CREATE TABLE IF NOT EXISTS eth.ancient_bodies (
  block_number BIGINT UNIQUE NOT NULL,
  body BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.ancient_bodies;