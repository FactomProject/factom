// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// ZeroHash is the string of all 00s
	ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	// RpcConfig sets the default target for the factomd and walletd API servers
	RpcConfig = &RPCConfig{
		FactomdServer: "localhost:8088",
		WalletServer:  "localhost:8089",
	}
)

// ChainIDFromFields computes a ChainID based on the binary External IDs of that
// Chain's First Entry.
func ChainIDFromFields(fields [][]byte) string {
	hs := sha256.New()
	for _, id := range fields {
		h := sha256.Sum256(id)
		hs.Write(h[:])
	}
	cid := hs.Sum(nil)
	return hex.EncodeToString(cid)
}

// ChainIDFromStrings computes the ChainID of a Chain Created with External IDs
// that would match the given string (in order).
func ChainIDFromStrings(fields []string) string {
	var bin [][]byte
	for _, str := range fields {
		bin = append(bin, []byte(str))
	}
	return ChainIDFromFields(bin)
}

// EntryCost calculates the cost in Entry Credits of adding an Entry to a Chain
// on the Factom protocol.
// The cost is the size of the Entry in Kilobytes excluding the Entry Header
// with any remainder being charged as a whole Kilobyte.
func EntryCost(e *Entry) (int8, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	}

	// caulculate the length exluding the header size 35 for Milestone 1
	l := len(p) - 35

	if l > 10240 {
		return 10, fmt.Errorf("Entry cannot be larger than 10KB")
	}

	// n is the capacity of the entry payment in KB
	n := int8(l / 1024)

	if r := l % 1024; r > 0 {
		n++
	}

	// The Entry Cost should never be less than one
	if n < 1 {
		n = 1
	}

	return n, nil
}

// FactoshiToFactoid converts a uint64 factoshi ammount into a fixed point
// number represented as a string
func FactoshiToFactoid(i uint64) string {
	d := i / 1e8
	r := i % 1e8
	ds := fmt.Sprintf("%d", d)
	rs := fmt.Sprintf("%08d", r)
	rs = strings.TrimRight(rs, "0")
	if len(rs) > 0 {
		ds = ds + "."
	}
	return fmt.Sprintf("%s%s", ds, rs)
}

// FactoidToFactoshi takes a Factoid amount as a string and returns the value in
// factoids
func FactoidToFactoshi(amt string) uint64 {
	valid := regexp.MustCompile(`^([0-9]+)?(\.[0-9]+)?$`)
	if !valid.MatchString(amt) {
		return 0
	}

	var total uint64 = 0

	dot := regexp.MustCompile(`\.`)
	pieces := dot.Split(amt, 2)
	whole, _ := strconv.Atoi(pieces[0])
	total += uint64(whole) * 1e8

	if len(pieces) > 1 {
		a := regexp.MustCompile(`(0*)([0-9]+)$`)

		as := a.FindStringSubmatch(pieces[1])
		part, _ := strconv.Atoi(as[0])
		power := len(as[1]) + len(as[2])
		total += uint64(part * 1e8 / int(math.Pow10(power)))
	}

	return total
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

// sha52 Sha512+Sha256 Hash; sha256(sha512(data)+data)
func sha52(data []byte) []byte {
	h1 := sha512.Sum512(data)
	h2 := sha256.Sum256(append(h1[:], data...))
	return h2[:]
}
