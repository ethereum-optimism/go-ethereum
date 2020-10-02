-- +goose Up
CREATE SCHEMA ethdb;

-- +goose Down
DROP SCHEMA ethdb;
