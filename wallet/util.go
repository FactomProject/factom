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
)

// seed address prefix
var seedPrefix = []byte{0x13, 0xdd}

// seedString returnes the string representation of a raw Wallet Seed or Next
// Wallet Seed.
func seedString(seed []byte) string {
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
