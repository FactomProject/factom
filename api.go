// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	ed "github.com/agl/ed25519"
	"golang.org/x/crypto/sha3"
)

var (
	server = "localhost:8088"
)

// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, key *[64]byte) error {
	type commit struct {
		CommitChainMsg string
	}

	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp
	buf.Write(milliTime())

	e := c.FirstEntry

	// 32 byte ChainID Hash
	if p, err := hex.DecodeString(c.ChainID); err != nil {
		return err
	} else {
		// double sha256 hash of ChainID
		buf.Write(shad(p))
	}

	// 32 byte Weld; sha256(sha256(EntryHash + ChainID))
	if cid, err := hex.DecodeString(c.ChainID); err != nil {
		return err
	} else {
		s := append(e.Hash(), cid...)
		buf.Write(shad(s))
	}

	// 32 byte Entry Hash of the First Entry
	buf.Write(e.Hash())

	// 1 byte number of Entry Credits to pay
	if d, err := ecCost(e); err != nil {
		return err
	} else {
		buf.WriteByte(byte(d + 10))
	}

	msg := buf.Bytes()

	// 32 byte Pubkey
	buf.Write(key[32:64])

	// 64 byte Signature of data from the Verstion to the Entry Credits
	buf.Write(ed.Sign(key, msg)[:])

	com := new(commit)
	com.CommitChainMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/commit-chain/", server),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// CommitEntry sends the signed Entry Hash and the Entry Credit public key to
// the factom network. Once the payment is verified and the network is commited
// to publishing the Entry it may be published with a call to RevealEntry.
func CommitEntry(e *Entry, key *[64]byte) error {
	type commit struct {
		CommitEntryMsg string
	}

	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp (truncated unix time)
	buf.Write(milliTime())

	// 32 byte Entry Hash
	buf.Write(e.Hash())

	// 1 byte number of entry credits to pay
	if c, err := ecCost(e); err != nil {
		return err
	} else {
		buf.WriteByte(byte(c))
	}

	// msg is the byte string before the pubkey and sig
	msg := buf.Bytes()

	// 32 byte public key
	buf.Write(key[32:64])

	// 64 byte signature
	buf.Write(ed.Sign(key, msg)[:])

	com := new(commit)
	com.CommitEntryMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/commit-entry/", server),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func RevealEntry(e *Entry) error {
	type reveal struct {
		Entry string
	}

	r := new(reveal)
	if p, err := e.MarshalBinary(); err != nil {
		return err
	} else {
		r.Entry = hex.EncodeToString(p)
	}

	j, err := json.Marshal(r)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/reveal-entry/", server),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

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

func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

func ecCost(e *Entry) (int8, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	}
	// n is the capacity of the entry payment in KB
	r := len(p) % 1024
	n := int8(len(p) / 1024)
	if r > 0 {
		n += 1
	}
	if n > 10 {
		return n, fmt.Errorf("Cannot make a payment for Entry larger than 10KB")
	}
	return n, nil
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
