// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
    "bytes"
    "testing"
)

func TestHexCompact(t *testing.T) {
	tests := []struct{ hex, compact []byte }{
		// empty keys, with and without terminator.
		{hex: []byte{}, compact: []byte{0x00}},
		{hex: []byte{16}, compact: []byte{0x20}},
		// odd length, no terminator
		{hex: []byte{1, 2, 3, 4, 5}, compact: []byte{0x11, 0x23, 0x45}},
		// even length, no terminator
		{hex: []byte{0, 1, 2, 3, 4, 5}, compact: []byte{0x00, 0x01, 0x23, 0x45}},
		// odd length, terminator
		{hex: []byte{15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x3f, 0x1c, 0xb8}},
		// even length, terminator
		{hex: []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x20, 0x0f, 0x1c, 0xb8}},
	}
	for _, test := range tests {
		if c := hexToCompact(test.hex); !bytes.Equal(c, test.compact) {
			t.Errorf("hexToCompact(%x) -> %x, want %x", test.hex, c, test.compact)
		}
		if h := compactToHex(test.compact); !bytes.Equal(h, test.hex) {
			t.Errorf("compactToHex(%x) -> %x, want %x", test.compact, h, test.hex)
		}
	}
}

func TestBinCompact(t *testing.T) {
	tests := []struct{ bin, compact []byte }{
		//empty keys, with and without terminator
		{bin: []byte{}, compact: []byte{0x40}}, // 0100 0000
		{bin: []byte{2}, compact: []byte{0xc0}}, // 1100 0000

		// length 1 with and without terminator
		{bin: []byte{1}, compact: []byte{0x38}}, // 0011 1000
		{bin: []byte{1, 2}, compact: []byte{0xb8}}, // 1011 1000

		// length 2 with and without terminator
		{bin: []byte{0,1}, compact: []byte{0x24}}, // 0010 0100
		{bin: []byte{0,1, 2}, compact: []byte{0xa4}}, // 1010 0100

		// length 3 with and without terminator
		{bin: []byte{1,0,1}, compact: []byte{0x1a}}, // 0001 1010
		{bin: []byte{1,0,1, 2}, compact: []byte{0x9a}}, // 1001 1010

		// length 4 with and without terminator
		{bin: []byte{1,0,1,0}, compact: []byte{0x0a}}, // 0000 1010
		{bin: []byte{1,0,1,0, 2}, compact: []byte{0x8a}}, // 1000 1010

		// length 5 with and without terminator
		{bin: []byte{1,0,1,0, 1}, compact: []byte{0x7a, 0x80}}, // 0111 1010 1000 0000
		{bin: []byte{1,0,1,0, 1, 2}, compact: []byte{0xfa, 0x80}}, // 1111 1010 1000 0000

		// length 6 with and without terminator
		{bin: []byte{1,0,1,0, 1,0}, compact: []byte{0x6a, 0x80}}, // 0110 1010 1000 0000
		{bin: []byte{1,0,1,0, 1,0, 2}, compact: []byte{0xea, 0x80}}, // 1110 1010 1000 0000

		// length 7 with and without terminator
		{bin: []byte{1,0,1,0, 1,0,1}, compact: []byte{0x5a, 0xa0}}, // 0101 1010 1010 0000
		{bin: []byte{1,0,1,0, 1,0,1, 2}, compact: []byte{0xda, 0xa0}}, // 1101 1010 1010 0000

		// length 8 with and without terminator
		{bin: []byte{1,0,1,0, 1,0,1,0}, compact: []byte{0x4a, 0xa0}}, // 0100 1010 1010 0000
		{bin: []byte{1,0,1,0, 1,0,1,0, 2}, compact: []byte{0xca, 0xa0}}, // 1100 1010 1010 0000
	}
	for _, test := range tests {
		if c := binaryToCompact(test.bin); !bytes.Equal(c, test.compact) {
			t.Errorf("binaryToCompact(%x) -> %x, want %x", test.bin, c, test.compact)
		}
		if h := compactToBinary(test.compact); !bytes.Equal(h, test.bin) {
			t.Errorf("compactToBinary(%x) -> %x, want %x", test.compact, h, test.bin)
		}
	}
}

func TestHexKeybytes(t *testing.T) {
	tests := []struct{ key, hexIn, hexOut []byte }{
		{key: []byte{}, hexIn: []byte{16}, hexOut: []byte{16}},
		{key: []byte{}, hexIn: []byte{}, hexOut: []byte{16}},
		{
			key:    []byte{0x12, 0x34, 0x56},
			hexIn:  []byte{1, 2, 3, 4, 5, 6, 16},
			hexOut: []byte{1, 2, 3, 4, 5, 6, 16},
		},
		{
			key:    []byte{0x12, 0x34, 0x5},
			hexIn:  []byte{1, 2, 3, 4, 0, 5, 16},
			hexOut: []byte{1, 2, 3, 4, 0, 5, 16},
		},
		{
			key:    []byte{0x12, 0x34, 0x56},
			hexIn:  []byte{1, 2, 3, 4, 5, 6},
			hexOut: []byte{1, 2, 3, 4, 5, 6, 16},
		},
	}
	for _, test := range tests {
		if h := keybytesToHex(test.key); !bytes.Equal(h, test.hexOut) {
			t.Errorf("keybytesToHex(%x) -> %x, want %x", test.key, h, test.hexOut)
		}
		if k := hexToKeybytes(test.hexIn); !bytes.Equal(k, test.key) {
			t.Errorf("hexToKeybytes(%x) -> %x, want %x", test.hexIn, k, test.key)
		}
	}
}

func TestBinaryKeybytes(t *testing.T) {
	tests := []struct{ key, binaryIn, binaryOut []byte }{
		{key: []byte{16}, binaryIn: []byte{2}, binaryOut: []byte{2}},
		{key: []byte{}, binaryIn: []byte{}, binaryOut: []byte{2}},
		{
			key:       []byte{1, 2, 3, 4, 5, 6, 16},
			binaryIn:  []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 0,1,1,0, 2},
			binaryOut: []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 0,1,1,0, 2},
		},
		{
			key:       []byte{1, 2, 3, 4, 0, 5, 16},
			binaryIn:  []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,0,0,0, 0,1,0,1, 2},
			binaryOut: []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,0,0,0, 0,1,0,1, 2},
		},
		{
			key:       []byte{1, 2, 3, 4, 5, 6},
			binaryIn:  []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 0,1,1,0},
			binaryOut: []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 0,1,1,0, 2},
		},
	}
	for _, test := range tests {
		if h := hexKeyBytesToBinary(test.key); !bytes.Equal(h, test.binaryOut) {
			t.Errorf("hexKeyBytesToBinary(%x) -> %b, want %b", test.key, h, test.binaryOut)
		}
		if k := binaryToHexKeyBytes(test.binaryIn); !bytes.Equal(k, test.key) {
			t.Errorf("binaryToHexKeyBytes(%b) -> %x, want %x", test.binaryIn, k, test.key)
		}
	}
}

func TestBinaryToHexKeyBytesWithPadding(t *testing.T) {
	tests := []struct{ binaryIn, hexOut []byte }{
		{
			binaryIn:  []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 1},
			hexOut:    []byte{1, 2, 3, 4, 5, 8},
		},
		{
			binaryIn:  []byte{0,0,0,1, 0,0,1,0, 0,0,1,1, 0,1,0,0, 0,1,0,1, 1, 2},
			hexOut:    []byte{1, 2, 3, 4, 5, 8, 16},
		},
	}
	for _, test := range tests {
		if k := binaryToHexKeyBytes(test.binaryIn); !bytes.Equal(k, test.hexOut) {
			t.Errorf("binaryToHexKeyBytes(%b) -> %x, want %x", test.binaryIn, k, test.hexOut)
		}
	}
}

func BenchmarkHexToCompact(b *testing.B) {
	testBytes := []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}
	for i := 0; i < b.N; i++ {
		hexToCompact(testBytes)
	}
}

func BenchmarkCompactToHex(b *testing.B) {
	testBytes := []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}
	for i := 0; i < b.N; i++ {
		compactToHex(testBytes)
	}
}

func BenchmarkKeybytesToHex(b *testing.B) {
	testBytes := []byte{7, 6, 6, 5, 7, 2, 6, 2, 16}
	for i := 0; i < b.N; i++ {
		keybytesToHex(testBytes)
	}
}

func BenchmarkHexToKeybytes(b *testing.B) {
	testBytes := []byte{7, 6, 6, 5, 7, 2, 6, 2, 16}
	for i := 0; i < b.N; i++ {
		hexToKeybytes(testBytes)
	}
}
