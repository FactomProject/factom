// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"time"
)

// FBlock represents a Factoid Block returned from factomd.
// Note: the FBlock api return does not use a "Header" field like the other
// block types do for some reason.
type FBlock struct {
	BodyMR          string           `json:"bodymr"`          // Merkle root of the Factoid transactions which accompany this block.
	PrevKeyMR       string           `json:"prevkeymr"`       // Key Merkle root of previous block.
	PrevLedgerKeyMR string           `json:"prevledgerkeymr"` // Sha3 of the previous Factoid Block
	ExchRate        int64            `json:"exchrate"`        // Factoshis per Entry Credit
	DBHeight        int64            `json:"dbheight"`        // Directory Block height
	Transactions    []*FBTransaction `json:"transactions"`    // The transactions inside the block

	ChainID     string `json:"chainid,omitempty"`
	KeyMR       string `json:"keymr,omitempty"`
	LedgerKeyMR string `json:"ledgerkeymr,omitempty"`
}

// FBTransactions represents a single and valid transaction contained inside
// of an FBlock.
// The data has been rearranged from the raw json response to make it easier to work with.
type FBTransaction struct {
	TxID        string                     `json:"txid"` // hex string
	BlockHeight int64                      `json:"blockheight"`
	Timestamp   time.Time                  `json:"timestamp"`
	Inputs      []SignedTransactionAddress `json:"inputs"`
	Outputs     []TransactionAddress       `json:"outputs"`
	ECOutputs   []TransactionAddress       `json:"outecs"`
}

// TransactionAddress holds the relevant data for either an input or an output.
// The amount is in either Factoshi (Input and Output) or EC (ECOutputs).
// The RCDHash is the SHA256 hash of the RCD.
// The address is the human readable address calculated from the RCDHash and type.
type TransactionAddress struct {
	Amount  uint64 `json:"amount"`
	RCDHash string `json:"address"` // hex string
	Address string `json:"useraddress"`
}

// SignedTransactionAddress contains a TransactionAddress along with the RCD and
// cryptographic signatures specified by the RCD.
type SignedTransactionAddress struct {
	TransactionAddress
	RCD        string   `json:"rcd"`        // hex string
	Signatures []string `json:"signatures"` // slice of hex strings
}

// these two types are just used internally as intermediary holding structs
// to transform the json response to the appropriate struct and back
type rawFBTransaction struct {
	TxID           string               `json:"txid"`
	BlockHeight    int64                `json:"blockheight"`
	MilliTimestamp int64                `json:"millitimestamp"`
	Inputs         []TransactionAddress `json:"inputs"`
	Outputs        []TransactionAddress `json:"outputs"`
	ECOutputs      []TransactionAddress `json:"outecs"`
	RCDs           []string             `json:"rcds"`
	SigBlocks      []rawSigBlock        `json:"sigblocks"`
}
type rawSigBlock struct {
	Signatures []string `json:"signatures"`
}

func (t *FBTransaction) MarshalJSON() ([]byte, error) {
	txReq := &rawFBTransaction{
		TxID:           t.TxID,
		BlockHeight:    t.BlockHeight,
		MilliTimestamp: t.Timestamp.UnixNano()/1e6 + (t.Timestamp.UnixNano()/1e3)%1e3,
		Outputs:        t.Outputs,
		ECOutputs:      t.ECOutputs,
		Inputs:         make([]TransactionAddress, 0, len(t.Inputs)),
		RCDs:           make([]string, 0, len(t.Inputs)),
		SigBlocks:      make([]rawSigBlock, 0, len(t.Inputs)),
	}

	for _, in := range t.Inputs {
		txReq.Inputs = append(txReq.Inputs, TransactionAddress{Amount: in.Amount, RCDHash: in.RCDHash, Address: in.Address})
		txReq.RCDs = append(txReq.RCDs, in.RCD)
		txReq.SigBlocks = append(txReq.SigBlocks, rawSigBlock{in.Signatures})
	}

	return json.Marshal(txReq)
}

func (t *FBTransaction) UnmarshalJSON(data []byte) error {
	txResp := new(rawFBTransaction)

	if err := json.Unmarshal(data, txResp); err != nil {
		return err
	}

	t.BlockHeight = txResp.BlockHeight
	// the bug in the nanosecond conversion is intentional to stay consistent with factomd
	t.Timestamp = time.Unix(txResp.MilliTimestamp/1e3, (txResp.MilliTimestamp%1e3)*1e3)
	t.Outputs = txResp.Outputs
	t.ECOutputs = txResp.ECOutputs
	t.TxID = txResp.TxID

	// catch decoding errors or malicious data
	if len(txResp.Inputs) != len(txResp.RCDs) || len(txResp.Inputs) != len(txResp.SigBlocks) {
		return fmt.Errorf("invalid signature counts")
	}

	for i := range txResp.Inputs {
		var sta SignedTransactionAddress
		sta.Amount = txResp.Inputs[i].Amount
		sta.RCDHash = txResp.Inputs[i].RCDHash
		sta.Address = txResp.Inputs[i].Address
		sta.RCD = txResp.RCDs[i]
		sta.Signatures = txResp.SigBlocks[i].Signatures

		t.Inputs = append(t.Inputs, sta)
	}

	return nil
}

func (t FBTransaction) String() string {
	var s string
	s += fmt.Sprintln("TxID:", t.TxID)
	s += fmt.Sprintln("BlockHeight:", t.BlockHeight)
	s += fmt.Sprintln("Timestamp:", t.Timestamp)

	if len(t.Inputs) > 0 {
		s += fmt.Sprintln("Inputs:")
		for _, in := range t.Inputs {
			s += fmt.Sprintln("   ", in.Address, in.Amount)
		}
	}

	if len(t.Outputs) > 0 {
		s += fmt.Sprintln("Outputs:")
		for _, out := range t.Outputs {
			s += fmt.Sprintln("   ", out.Address, out.Amount)
		}
	}

	if len(t.ECOutputs) > 0 {
		s += fmt.Sprintln("ECOutputs:")
		for _, ec := range t.ECOutputs {
			s += fmt.Sprintln("   ", ec.Address, ec.Amount)
		}
	}

	return s
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
		s += fmt.Sprintln(t.String())
	}
	s += fmt.Sprintln("}")

	return s
}

// GetFBlock requests a specified Factoid Block from factomd by its keymr
func GetFBlock(keymr string) (fblock *FBlock, err error) {
	params := keyMRRequest{KeyMR: keymr, NoRaw: true}
	req := NewJSON2Request("factoid-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	// Create temporary struct to unmarshal json object
	wrap := new(struct {
		FBlock *FBlock `json:"fblock"`
	})

	if err = json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return
	}

	return wrap.FBlock, nil
}

// GetFBlockByHeight requests a specified Factoid Block from factomd by its height
func GetFBlockByHeight(height int64) (fblock *FBlock, err error) {
	params := heightRequest{Height: height, NoRaw: true}
	req := NewJSON2Request("fblock-by-height", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	wrap := new(struct {
		FBlock *FBlock `json:"fblock"`
	})
	if err = json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return
	}

	return wrap.FBlock, nil
}
