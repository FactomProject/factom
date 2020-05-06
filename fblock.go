// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// FBlock represents a Factoid Block returned from factomd.
// Note: the FBlock api return does not use a "Header" field like the other
// block types do for some reason.
type FBlock struct {
	BodyMR          string            `json:"bodymr"`          // Merkle root of the Factoid transactions which accompany this block.
	PrevKeyMR       string            `json:"prevkeymr"`       // Key Merkle root of previous block.
	PrevLedgerKeyMR string            `json:"prevledgerkeymr"` // Sha3 of the previous Factoid Block
	ExchRate        int64             `json:"exchrate"`        // Factoshis per Entry Credit
	DBHeight        int64             `json:"dbheight"`        // Directory Block height
	Transactions    []json.RawMessage `json:"transactions"`

	ChainID     string `json:"chainid,omitempty"`
	KeyMR       string `json:"keymr,omitempty"`
	LedgerKeyMR string `json:"ledgerkeymr,omitempty"`
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
		s += fmt.Sprintln(string(t))
	}
	s += fmt.Sprintln("}")

	return s
}

func getFBlock(keymr string, noraw bool) (fblock *FBlock, raw []byte, err error) {
	params := keyMRRequest{KeyMR: keymr, NoRaw: noraw}
	req := NewJSON2Request("factoid-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	// Create temporary struct to unmarshal json object
	wrap := new(struct {
		FBlock  *FBlock `json:"fblock"`
		RawData string  `json:"rawdata"`
	})

	if err = json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return
	}

	raw, err = hex.DecodeString(wrap.RawData)
	if err != nil {
		return
	}

	return wrap.FBlock, raw, nil
}

// GetFBlock requests a specified Factoid Block from factomd. It returns the
// FBlock struct, the raw binary FBlock, and an error if present.
func GetFBlock(keymr string) (fblock *FBlock, raw []byte, err error) {
	return getFBlock(keymr, false)
}

// GetSimpleFBlock requests a specified Factoid Block from factomd. It returns the
// FBlock struct, and an error if present.
func GetSimpleFBlock(keymr string) (fblock *FBlock, err error) {
	fblock, _, err = getFBlock(keymr, true)
	return
}

func getFBlockByHeight(height int64, noraw bool) (fblock *FBlock, raw []byte, err error) {
	params := heightRequest{Height: height, NoRaw: noraw}
	req := NewJSON2Request("fblock-by-height", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	wrap := new(struct {
		FBlock  *FBlock `json:"fblock"`
		RawData string  `json:"rawdata"`
	})
	if err = json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return
	}

	raw, err = hex.DecodeString(wrap.RawData)
	if err != nil {
		return
	}

	return wrap.FBlock, raw, nil
}

// GetFBlockByHeight requests a specified Factoid Block from factomd, returning
// the FBlock struct, the raw binary FBlock, and an error if present.
func GetFBlockByHeight(height int64) (fblock *FBlock, raw []byte, err error) {
	return getFBlockByHeight(height, false)
}

// GetSimpleFBlockByHeight requests a specified Factoid Block from factomd, returning
// the FBlock struct, the raw binary FBlock, and an error if present.
func GetSimpleFBlockByHeight(height int64) (fblock *FBlock, err error) {
	fblock, _, err = getFBlockByHeight(height, true)
	return
}
