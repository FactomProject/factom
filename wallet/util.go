// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"crypto/sha256"

	"github.com/FactomProject/btcutil/base58"
)

const (
	SeedLength = 64
	Version    = "0.2.1.0"
	ApiVersion = "2.0"
)

// seed address prefix
var seedPrefix = []byte{0x13, 0xdd}

// SeedString returnes the string representation of a raw Wallet Seed or Next
// Wallet Seed.
func SeedString(seed []byte) string {
	if len(seed) != SeedLength {
		return ""
	}

	buf := new(bytes.Buffer)

	// 2 byte Seed Address Prefix
	buf.Write(seedPrefix)

	// 64 byte Seed
	buf.Write(seed)

	// 4 byte Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

// shad Double Sha256 Hash; sha256(sha256(data))
func shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// newCounter is used to generate the ID field for the JSON2Request
func newCounter() func() int {
	count := 0
	return func() int {
		count += 1
		return count
	}
}

var APICounter = newCounter()
