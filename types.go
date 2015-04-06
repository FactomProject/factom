// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
	"time"
)

type Chain struct {
	Blocks     []EBlock
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
