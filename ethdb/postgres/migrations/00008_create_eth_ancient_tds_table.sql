-- +goose Up
CREATE TABLE IF NOT EXISTS eth.ancient_tds (
  block_number BIGINT UNIQUE NOT NULL,
  td BYTEA NOT NULL
);

-- +goose Down
DROP TABLE eth.ancient_tds;