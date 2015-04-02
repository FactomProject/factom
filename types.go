// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"time"
)

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
	EBEntries []struct {
		TimeStamp int64
		Hash      string
	}
}

func (e EBlock) Time() time.Time {
	return time.Unix(e.Header.TimeStamp, 0)
}

type Entry struct {
	ChainID string
	ExtIDs  []string
	Data    string
}

// TODO
type Chain struct {
}
