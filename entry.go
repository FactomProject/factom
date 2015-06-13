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
	"io/ioutil"
	"net/http"
	"time"

	ed "github.com/agl/ed25519"
)

type Entry struct {
	ChainID string
	ExtIDs  []string
	Content string
}

func NewEntry() *Entry {
	e := new(Entry)
	
	return e
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
	if c, err := entryCost(e); err != nil {
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

func GetEntry(hash string) (*Entry, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/entry-by-hash/%s", server, hash))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	e := new(Entry)
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}
	
	return e, nil
}

func (e *Entry) Hash() []byte {
	a, err := e.MarshalBinary()
	if err != nil {
		return make([]byte, 32)
	}
	return sha23(a)
}

func (e *Entry) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	c, err := hex.DecodeString(e.Content)
	if err != nil {
		return buf.Bytes(), err
	}
	x, err := e.MarshalExtIDsBinary()
	if err != nil {
		return buf.Bytes(), err
	}

	// Header

	// 1 byte Version
	buf.Write([]byte{0})

	// 32 byte chainid
	if p, err := hex.DecodeString(e.ChainID); err != nil {
		return buf.Bytes(), err
	} else {
		buf.Write(p)
	}

	// 2 byte size of extids
	if err := binary.Write(buf, binary.BigEndian, int16(len(x))); err != nil {
		return buf.Bytes(), err
	}

	// Payload

	// extids
	buf.Write(x)

	// data
	buf.Write(c)

	return buf.Bytes(), nil
}

func (e *Entry) MarshalExtIDsBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, v := range e.ExtIDs {
		p, err := hex.DecodeString(v)
		if err != nil {
			return buf.Bytes(), err
		}
		// 2 byte length of extid
		binary.Write(buf, binary.BigEndian, int16(len(p)))
		// extid
		buf.Write(p)
	}

	return buf.Bytes(), nil
}

func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

func entryCost(e *Entry) (int8, error) {
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
