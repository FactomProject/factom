// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// EBlock is an Entry Block from the Factom Network. An Entry Block contains a
// series of Entries all belonging to the same Chain on Factom from a given 10
// minute period. All of the Entry Blocks from a given period are collected into
// a Merkel Tree the root of which is the Factom Directory Block.
type EBlock struct {
	Header struct {
		// TODO: rename BlockSequenceNumber to EBSequence
		BlockSequenceNumber int64  `json:"blocksequencenumber"`
		ChainID             string `json:"chainid"`
		PrevKeyMR           string `json:"prevkeymr"`
		Timestamp           int64  `json:"timestamp"`
		DBHeight            int64  `json:"dbheight"`
	} `json:"header"`
	EntryList []EBEntry `json:"entrylist"`
}

// EBEntry is a member of the Entry Block representing a Factom Entry. The
// EBEntry has the hash of the Factom Entry and the time when the entry was
// added.
//
// The consensus algorithm does NOT garuntee that the cannonical order of the
// EBEntries is the same order that they were recieved for multiple Entries that
// are made to the network during a single minute.
type EBEntry struct {
	EntryHash string `json:"entryhash"`
	Timestamp int64  `json:"timestamp"`
}

func (e *EBlock) String() string {
	var s string
	s += fmt.Sprintln("BlockSequenceNumber:", e.Header.BlockSequenceNumber)
	s += fmt.Sprintln("ChainID:", e.Header.ChainID)
	s += fmt.Sprintln("PrevKeyMR:", e.Header.PrevKeyMR)
	s += fmt.Sprintln("Timestamp:", e.Header.Timestamp)
	s += fmt.Sprintln("DBHeight:", e.Header.DBHeight)
	for _, v := range e.EntryList {
		s += fmt.Sprintln("EBEntry {")
		s += fmt.Sprintln("	Timestamp", v.Timestamp)
		s += fmt.Sprintln("	EntryHash", v.EntryHash)
		s += fmt.Sprintln("}")
	}
	return s
}

// GetEBlock requests an Entry Block from factomd by its Key Merkle Root
func GetEBlock(keymr string) (*EBlock, error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("entry-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(EBlock)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}

// GetAllEBlockEntries requests every Entry from a given Entry Block
func GetAllEBlockEntries(keymr string) ([]*Entry, error) {
	es := make([]*Entry, 0)

	eb, err := GetEBlock(keymr)
	if err != nil {
		return es, err
	}

	for _, v := range eb.EntryList {
		e, err := GetEntry(v.EntryHash)
		if err != nil {
			return es, err
		}
		es = append(es, e)
	}

	return es, nil
}
