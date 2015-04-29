// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"
	
	"golang.org/x/crypto/sha3"
)

type Chain struct {
	ChainID    string
	FirstEntry Entry
	Name       []string
}

func (c *Chain) Names() []string {
	names := make([]string, 0)
	for _, v := range c.Name {
		p, err := hex.DecodeString(v)
		if err != nil {
			return nil
		}
		names = append(names, string(p))
	}
	return names
}

type DBInfo struct {
	BTCBlockHash   string
	BTCBlockHeight int64
	BTCTxHash      string
	BTCTxOffset    int64
	DBHash         string
	DBMerkleRoot   string
}

type DBlock struct {
	DBHash string
	Header struct {
		BlockID       int
		EntryCount    int
		MerkleRoot    string
		PrevBlockHash string
		TimeStamp     int64
	}
	DBEntries []struct {
		ChainID    string
		MerkleRoot string
	}
}

func (d DBlock) Time() time.Time {
	return time.Unix(d.Header.TimeStamp, 0)
}

type EBlock struct {
	Header struct {
		BlockID       int
		PrevBlockHash string
		TimeStamp     int64
	}
	EBEntries []EBEntry
}

func (e EBlock) Time() time.Time {
	return time.Unix(e.Header.TimeStamp, 0)
}

type EBEntry struct {
	TimeStamp int64
	Hash      string
}

func (e EBEntry) Time() time.Time {
	return time.Unix(e.TimeStamp, 0)
}

type Entry struct {
	ChainID string
	ExtIDs  []string
	Data    string
}

func (e *Entry) Hash() [32]byte {
	a, err := e.MarshalBinary()
	if err != nil {
		return [32]byte{byte(0)}
	}
	fmt.Println("a:", hex.EncodeToString(a[:]))
	b := sha3.Sum256(a)
	fmt.Println("b:", hex.EncodeToString(b[:]))
	c := append(a, b[:]...)
	
	return sha256.Sum256(c)
}

func (e *Entry) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	d, err := hex.DecodeString(e.Data)
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
	
	// 2 byte payload size
	if err := binary.Write(buf, binary.BigEndian, int16(len(x) + len(d)));
		err != nil {
		return buf.Bytes(), err
	}
	
	// Payload
	// extids
	buf.Write(x)
	
	// content
	buf.Write(d)
	
	return buf.Bytes(), nil
}

func (e *Entry) MarshalExtIDsBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	
	for _, v := range e.ExtIDs {
		p, err := hex.DecodeString(v)
		if err != nil {
			return buf.Bytes(), err
		}
		binary.Write(buf, binary.BigEndian, int16(len(p)))
		buf.Write(p)
	}
	
	return buf.Bytes(), nil
}
