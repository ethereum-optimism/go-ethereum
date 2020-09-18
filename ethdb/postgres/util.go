// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package postgres

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-ds-help"
	_ "github.com/lib/pq" //postgres driver
	"github.com/multiformats/go-multihash"
)

// MultihashKeyFromKeccak256 converts keccak256 hash bytes into a blockstore-prefixed multihash db key string
func MultihashKeyFromKeccak256(h []byte) (string, error) {
	mh, err := multihash.Encode(h, multihash.KECCAK_256)
	if err != nil {
		return "", err
	}
	dbKey := dshelp.MultihashToDsKey(mh)
	return blockstore.BlockPrefix.String() + dbKey.String(), nil
}

// DatastoreKeyFromGethKey returns the public.blocks key from the provided geth key
// It also returns the key's prefix, if it has one
func DatastoreKeyFromGethKey(h []byte) (string, []byte, error) {
	keyType, keyComponents := ResolveKeyType(h)
	switch keyType {
	case Keccak:
		mhKey, err := MultihashKeyFromKeccak256(h)
		return mhKey, nil, err
	case Header:
		mhKey, err := MultihashKeyFromKeccak256(keyComponents[1])
		return mhKey, keyComponents[0], err
	case Preimage:
		mhKey, err := MultihashKeyFromKeccak256(keyComponents[1])
		return mhKey, keyComponents[0], err
	case Prefixed, Suffixed:
		// This data is not mapped by hash => content by geth, store it using the prefixed/suffixed key directly
		// I.e. the public.blocks datastore key == the hex representation of the geth key
		// Alternatively, decompose the data and derive the hash
		return common.Bytes2Hex(h), keyComponents[0], nil
	case Static:
		return common.Bytes2Hex(h), nil, nil
	default:
		return "", nil, fmt.Errorf("invalid formatting of database key: %x", h)
	}
}
