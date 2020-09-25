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
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/rawdb"

	_ "github.com/lib/pq" //postgres driver
)

// ResolveKeyPrefix returns the key and its prefix, if it has one
func ResolveKeyPrefix(key []byte) ([]byte, []byte, error) {
	sk := bytes.Split(key, rawdb.PrefixDelineation)
	switch l := len(sk); {
	case l == 1:
		return key, nil, nil
	case l == 2:
		return key, sk[0], nil
	default:
		return nil, nil, fmt.Errorf("unexpected number of key components: %d", l)
	}
}
