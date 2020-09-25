-- +goose Up
CREATE TABLE IF NOT EXISTS eth.kvstore (
  eth_key BYTEA UNIQUE NOT NULL,
  eth_data BYTEA NOT NULL,
  prefix BYTEA
);

CREATE INDEX prefix_index ON eth.kvstore USING btree (prefix);

-- +goose Down
DROP TABLE eth.kvstore;
DROP INDEX eth.prefix_index;
