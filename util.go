// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"crypto/sha256"
	"os"

	ed "github.com/agl/ed25519"
	"golang.org/x/crypto/sha3"
)

var (
	server = "localhost:8088"
    serverFct = "localhost:8089"
)

func NewECKey() *[64]byte {
	rand, err := os.Open("/dev/random")
	if err != nil {
		return &[64]byte{byte(0)}
	}

	// private key is [32]byte private section + [32]byte public key
	_, priv, err := ed.GenerateKey(rand)
	if err != nil {
		return &[64]byte{byte(0)}
	}
	return priv
}

func SetServer(s string) {
	server = s
}

// shad Double Sha256 Hash; sha256(sha256(data))
func shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// sha23 combination sha256 and sha3 Hash; sha256(data + sha3(data))
func sha23(data []byte) []byte {
	h1 := sha3.Sum256(data)
	h2 := sha256.Sum256(append(data, h1[:]...))
	return h2[:]
}
