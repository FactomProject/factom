// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"time"
)

type DBlock struct {
	Header struct {
		Version       int
		TimeStamp     int64
		BatchFlag     int
		EntryCount    int
		BlockID       int
		PrevBlockHash string
		MerkleRoot    string
	}
	DBEntries []struct {
		MerkleRoot string
		ChainID    string
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
		timeStamp int64
		hash      string
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
