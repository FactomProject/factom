// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

type Chain struct {
	ChainID    string
	FirstEntry *Entry
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
