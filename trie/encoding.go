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

import "math"

// Trie keys are dealt with in three distinct encodings:
//
// KEYBYTES encoding contains the actual key and nothing else. This encoding is the
// input to most API functions.
//
// HEX encoding contains one byte for each nibble of the key and an optional trailing
// 'terminator' byte of value 0x10 which indicates whether or not the node at the key
// contains a value. Hex key encoding is used for nodes loaded in memory because it's
// convenient to access.
//
// COMPACT encoding is defined by the Ethereum Yellow Paper (it's called "hex prefix
// encoding" there) and contains the bytes of the key and a flag. The high nibble of the
// first byte contains the flag; the lowest bit encoding the oddness of the length and
// the second-lowest encoding whether the node at the key is a value node. The low nibble
// of the first byte is zero in the case of an even number of nibbles and the first nibble
// in the case of an odd number. All remaining nibbles (now an even number) fit properly
// into the remaining bytes. Compact encoding is used for nodes stored on disk.

// Proposal:
// HEADER NIBBLE:
// first bit: 1 if should be terminated / 0 if not,
// bits 2-4: the number of unused (least significant) bits from last byte = [8 - ((4 (header nibble) + message len) % 8)] % 8

// 1 1 0 1 1 2(terminator)
// first bit = 1 (terminator present)
// bits 2-4 = 4 + 5(msg len w/o terminator) % 8 = 001
// first nibble: 1001
// entire message = 1001 1101 1[000 0000], where [padding]
// encoded: 9d8


func binaryToCompact(bin []byte) []byte {
	currentByte := uint8(0)
	keyLength := len(bin)
	if hasBinaryTerminator(bin) {
		bin = bin[:len(bin)-1]
		currentByte = 1 << 7
		keyLength--
	}

	lastByteUnusedBits := uint8((8 - ((4 + keyLength)) % 8) % 8)
	currentByte += lastByteUnusedBits << 4

	returnLength := (keyLength + 4 + int(lastByteUnusedBits)) / 8
	returnBytes := make([]byte, returnLength)
	returnIndex := 0
	for i := 0; i < len(bin); i++ {
		bitPosition := (4 + i) % 8
		if bitPosition == 0 {
			returnBytes[returnIndex] = currentByte
			currentByte = uint8(0)
			returnIndex++
		}

		currentByte += (1 & bin[i]) << (7-bitPosition)
	}
	returnBytes[returnIndex] = currentByte

	return returnBytes
}

func compactToBinary(compact []byte) []byte {
	addTerminator := compact[0] >> 7
	lastByteUnusedBits := (compact[0] << 1) >> 5

	returnLength := len(compact) * 8 - 4        // length - header nibble
	returnLength += int(addTerminator)          // terminator byte
	returnLength -= int(lastByteUnusedBits)     // extra padding bits

	returnBytes := make([]byte, returnLength)

	returnIndex := 0
	byteIndex := 0
	bitPosition := 4
	currentByte := compact[byteIndex]
	for ; returnIndex < returnLength - int(addTerminator); bitPosition++ {
		shift := 7 - (bitPosition % 8)
		if shift == 7 {
			byteIndex++
			currentByte = compact[byteIndex]
		}
		returnBytes[returnIndex] = uint8((currentByte & (1 << shift)) >> shift)
		returnIndex++
	}

	if addTerminator > 0 {
		returnBytes[returnLength -1] = binTerminator
	}

	return returnBytes
}


func hexToCompact(hex []byte) []byte {
	terminator := byte(0)
	if hasTerm(hex) {
		terminator = 1
		hex = hex[:len(hex)-1]
	}
	buf := make([]byte, len(hex)/2+1)
	buf[0] = terminator << 5 // the flag byte
	if len(hex)&1 == 1 {
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= hex[0] // first nibble is contained in the first byte
		hex = hex[1:]
	}
	decodeNibbles(hex, buf[1:])
	return buf
}

func compactToHex(compact []byte) []byte {
	if len(compact) == 0 {
		return compact
	}
	base := keybytesToHex(compact)
	// delete terminator flag
	if base[0] < 2 {
		base = base[:len(base)-1]
	}
	// apply odd flag
	chop := 2 - base[0]&1
	return base[chop:]
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	//println(fmt.Sprintf("key bytes to hex. Start: %x, End: %x", str, nibbles))
	return nibbles
}

// hexToKeybytes turns hex nibbles into key bytes.
// This can only be used for keys of even length.
func hexToKeybytes(hex []byte) []byte {
	if hasTerm(hex) {
		hex = hex[:len(hex)-1]
	}
	if len(hex)&1 != 0 {
		panic("can't convert hex key of odd length")
	}
	key := make([]byte, len(hex)/2)
	decodeNibbles(hex, key)
	return key
}

// Converts a []byte key with hex granularity (4 unused bits + 4 used bits)
// to a []byte key with bit granularity (7 unused bits + 1 used bit).
func hexKeyBytesToBinary(hexKey []byte) (bitKey []byte) {
	length := len(hexKey) * 4
	if hasTerm(hexKey) {
		length -= 3
	} else {
		// add terminator byte
		length += 1
	}
	if hexKey == nil || len(hexKey) == 0 {
		ret := make([]byte, 1)
		ret[0] = binTerminator
		return ret
	}
	bitKey = make([]byte, length)

	for bite := 0; bite < length / 4; bite++ {
		for bit := 0; bit < 4; bit++ {
			// set right-most bit to 0 or 1 based whether the bit index of hexKey[bite] is set
			bitKey[(bite * 4) + bit] = (hexKey[bite] & uint8(8 >> bit)) >> (3-bit)
		}
	}

	if length % 2 != 0 {
		bitKey[len(bitKey) -1] = binTerminator
	}
	// println(fmt.Sprintf("hex key bytes to bin. Start: %x, End: %x, length: %d", hexKey, bitKey, length))

	return bitKey
}

// Converts a []byte key with bit granularity (max of one bit per byte set) to the half-packed bytes
// representing the hex-encoded version of the key.
func binaryToHexKeyBytes(bitKey []byte) (hexKey []byte) {
	addTerminator := 0
	if hasBinaryTerminator(bitKey) {
		bitKey = bitKey[:len(bitKey)-1]
		addTerminator = 1
	}
	if bitKey == nil || len(bitKey) == 0 {
		if addTerminator > 0 {
			ret := make([]byte, 1)
			ret[0] = terminator
			return ret
		}
		return make([]byte, 0)
	}

	paddingBits := len(bitKey) % 4
	//if len(bitKey) % 4 != 0 {
	//	panic(fmt.Sprintf("can't convert binary key of length %d to hex. Key: %x", len(bitKey), bitKey))
	//}
	hexKey = make([]byte, len(bitKey) / 4 + paddingBits + 1 * addTerminator)

	nibbleInt := uint8(0)
	for bit := 0; bit < len(bitKey); bit++ {
		nibbleBit := bit % 4
		if nibbleBit == 0 && bit != 0 {
			hexKey[(bit / 4) - 1] = nibbleInt
			nibbleInt = 0
		}
		nibbleInt += uint8(math.Pow(2, float64(3-nibbleBit))) * bitKey[bit]
	}

	hexKey[len(hexKey) -1 -addTerminator] = nibbleInt
	if addTerminator > 0 {
		hexKey[len(hexKey) - 1] = terminator
	}

	return hexKey
}

func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}

// prefixLen returns the length of the common prefix of a and b.
func prefixLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

const terminator = 16
// hasTerm returns whether a hex key has the terminator flag.
func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == terminator
}

const binTerminator = 2
// hasTerm returns whether a hex key has the terminator flag.
func hasBinaryTerminator(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == binTerminator
}
