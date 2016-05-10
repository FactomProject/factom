// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
)

type DBlock struct {
	DBHash string
	Header struct {
		PrevBlockKeyMR string
		Timestamp      uint64
		SequenceNumber int
	}
	EntryBlockList []struct {
		ChainID string
		KeyMR   string
	}
}

func (d *DBlock) String() string {
	var s string
	s += fmt.Sprintln("PrevBlockKeyMR:", d.Header.PrevBlockKeyMR)
	s += fmt.Sprintln("Timestamp:", d.Header.Timestamp)
	s += fmt.Sprintln("SequenceNumber:", d.Header.SequenceNumber)
	for _, v := range d.EntryBlockList {
		s += fmt.Sprintln("EntryBlock {")
		s += fmt.Sprintln("	ChainID", v.ChainID)
		s += fmt.Sprintln("	KeyMR", v.KeyMR)
		s += fmt.Sprintln("}")
	}
	return s
}

type DBHead struct {
	KeyMR string
}

type DirectoryBlockHeightResponse struct {
	Height int64
}
