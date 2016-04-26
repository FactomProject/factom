// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
)

type EBlock struct {
	Header struct {
		BlockSequenceNumber int
		ChainID             string
		PrevKeyMR           string
		Timestamp           uint64
	}
	EntryList []EBEntry
}

type EBEntry struct {
	Timestamp int64
	EntryHash string
}

func (e *EBlock) String() string {
	var s string
	s += fmt.Sprintln("BlockSequenceNumber:", e.Header.BlockSequenceNumber)
	s += fmt.Sprintln("ChainID:", e.Header.ChainID)
	s += fmt.Sprintln("PrevKeyMR:", e.Header.PrevKeyMR)
	s += fmt.Sprintln("Timestamp:", e.Header.Timestamp)
	for _, v := range e.EntryList {
		s += fmt.Sprintln("EBEntry {")
		s += fmt.Sprintln("	Timestamp", v.Timestamp)
		s += fmt.Sprintln("	EntryHash", v.EntryHash)
		s += fmt.Sprintln("}")
	}
	return s
}
