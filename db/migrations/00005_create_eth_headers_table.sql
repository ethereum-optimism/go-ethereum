-- +goose Up
CREATE TABLE IF NOT EXISTS eth.headers (
  header_key BYTEA PRIMARY KEY,
  header BYTEA NOT NULL,
  height BIGINT NOT NULL
);

CREATE INDEX header_height_index ON eth.headers USING brin (height);

-- +goose Down
DROP TABLE eth.headers;
DROP INDEX eth.header_height_index;