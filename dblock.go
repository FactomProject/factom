// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"
)

type DBlock struct {
	DBHash string `json:"dbhash"`
	Header struct {
		PrevBlockKeyMR string `json:"prevblockkeymr"`
		SequenceNumber int64  `json:"sequencenumber"`
		Timestamp      int64  `json:"timestamp"`
	} `json:"header"`
	EntryBlockList []struct {
		ChainID string `json:"chainid"`
		KeyMR   string `json:"keymr"`
	} `json:"entryblocklist"`
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
	KeyMR string `json:"keymr"`
}
