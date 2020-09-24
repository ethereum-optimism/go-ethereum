-- +goose Up
CREATE TABLE IF NOT EXISTS eth.key_preimages (
  eth_key BYTEA UNIQUE NOT NULL,
  prefix BYTEA,
  ipfs_key TEXT NOT NULL REFERENCES public.blocks (key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

-- +goose Down
DROP TABLE eth.key_preimages;
