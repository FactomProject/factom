// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

type Chain struct {
	ChainID    string
	FirstEntry Entry
}

type DBlock struct {
	DBHash string
	Header struct {
		PrevBlockKeyMR string
		TimeStamp      uint64
		SequenceNumber int
	}
	DBEntries []struct {
		ChainID string
		KeyMR   string
	}
}

type EBlock struct {
	Header struct {
		BlockSequenceNumber int
		ChainID             string
		PrevKeyMR           string
		TimeStamp           uint64
	}
	EBEntries []EBEntry
}

type EBEntry struct {
	TimeStamp int64
	EntryHash string
}

type Entry struct {
	ChainID string
	ExtIDs  []string
	Content string
}

func (e *Entry) Hash() [32]byte {
	a, err := e.MarshalBinary()
	if err != nil {
		return [32]byte{byte(0)}
	}
	b := sha3.Sum256(a)
	c := append(a, b[:]...)
	return sha256.Sum256(c)
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
