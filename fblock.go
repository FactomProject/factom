// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// FBlock represents a Factoid Block returned from factomd.
// Note: the FBlock api return does not use a "Header" field like the other
// block types do for some reason.
type FBlock struct {
	BodyMR          string `json:"bodymr"`          // Merkle root of the Factoid transactions which accompany this block.
	PrevKeyMR       string `json:"prevkeymr"`       // Key Merkle root of previous block.
	PrevLedgerKeyMR string `json:"prevledgerkeymr"` // Sha3 of the previous Factoid Block
	ExchRate        int64  `json:"exchrate"`        // Factoshis per Entry Credit
	DBHeight        int64  `json:"dbheight"`        // Directory Block height

	Transactions []Transaction `json:"transactions"`
}

func (f *FBlock) String() string {
	var s string

	s += fmt.Sprintln("BodyMR:", f.BodyMR)
	s += fmt.Sprintln("PrevKeyMR:", f.PrevKeyMR)
	s += fmt.Sprintln("PrevLedgerKeyMR:", f.PrevLedgerKeyMR)
	s += fmt.Sprintln("ExchRate:", f.ExchRate)
	s += fmt.Sprintln("DBHeight:", f.DBHeight)

	s += fmt.Sprintln("Transactions {")
	for _, t := range f.Transactions {
		s += fmt.Sprintln(t)
	}
	s += fmt.Sprintln("}")

	return s
}

// GetFblock requests a specified Factoid Block from factomd.
func GetFBlock(keymr string) (*FBlock, error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("factoid-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	// Create temporary struct to unmarshal json object
	f := new(struct {
		FBlock *FBlock `json:"fblock"`
	})

	if err := json.Unmarshal(resp.JSONResult(), f); err != nil {
		return nil, err
	}

	return f.FBlock, nil
}
