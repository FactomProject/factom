// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"time"
)

const (
	ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	server    = "localhost:8088"
	serverFct = "localhost:8089"
)

// SetServer sets the gloabal target for the factomd server
func SetServer(s string) {
	server = s
}

// SetWallet sets the global target for the fctwallet server
func SetWallet(s string) {
	serverFct = s
}

// Server() returns the global server string for debugging
func Server() string {
	return server
}

// milliTime returns a 6 byte slice representing the unix time in milliseconds
func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

// shad Double Sha256 Hash; sha256(sha256(data))
func shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// sha52
func sha52(data []byte) []byte {
	h1 := sha512.Sum512(data)
	h2 := sha256.Sum256(append(h1[:], data...))
	return h2[:]
}
