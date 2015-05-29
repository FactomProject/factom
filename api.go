// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	ed "github.com/agl/ed25519"
)

var (
	server = "localhost:8088"
)

/* TODO finish CommitChain
// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, key *[64]byte) error {
	buf := new(bytes.Buffer)
	
	// 1 byte version
	buf.Write([]byte{0})
	
	// 6 byte milliTimestamp (truncated unix time)
	m := milliTime()
	buf.Write(m)

	// 32 byte ChainID Hash
	if p, err := hex.DecodeString(e.ChainID); err != nil {
		return err
	} else {
		// double sha256 hash of ChainID
		h1 := sha256.Sum256(p)
		h2 := sha256.Sum256(h1[:])
		buf.Write(h2[:])
	}
	
	// 32 byte Hash of the Entry Hash + ChainID
	
	// 32 byte Entry Hash of the First Entry
	buf.Write(c.FirstEntry.Hash()[:])
	
	// 1 byte number of Entry Credits to pay
	if d, err := ecCost(c); err != nil {
		return err
	} else {
		buf.WriteByte(d)
	}
	
	msg := buf.Bytes()
	
	// 32 byte Pubkey
	buf.Write(key[32:64])
	
	// 64 byte Signature of data from the Verstion to the Entry Credits
	buf.Write(ed.Sign(key, msg))
}
*/

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
	h := e.Hash()
	buf.Write(h[:])
	
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
	
	c := new(commit)
	c.CommitEntryMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(c)
	if err != nil {
		return err
	}
	
	api := fmt.Sprintf("http://%s/v1/commit-entry/", server)
	resp, err := http.Post(api, "application/json", bytes.NewBuffer(j))
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
	
	api := fmt.Sprintf("http://%s/v1/reveal-entry/", server)
	resp, err := http.Post(api, "application/json", bytes.NewBuffer(j))
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
	r := len(p) % 1000
	n := int8(len(p) / 1000)
	if r > 0 {
		n += 1
	}
	if n > 10 {
		return n, fmt.Errorf("Entry larger than 10KB")
	}
	return n, nil
}
