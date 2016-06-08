// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

func GetReceipt(hash string) (*Receipt, error) {
	param := HashRequest{Hash: hash}
	req := NewJSON2Request("receipt", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	rec := new(ReceiptResponse)
	if err := json.Unmarshal(resp.Result, rec); err != nil {
		return nil, err
	}

	return rec.Receipt, nil
}

type ReceiptResponse struct {
	Receipt *Receipt `json:"receipt"`
}

type Receipt struct {
	Entry                  *JSON         `json:"entry,omitempty"`
	MerkleBranch           []*MerkleNode `json:"merklebranch,omitempty"`
	EntryBlockKeyMR        string        `json:"entryblockkeymr,omitempty"`
	DirectoryBlockKeyMR    string        `json:"directoryblockkeymr,omitempty"`
	BitcoinTransactionHash string        `json:"bitcointransactionhash,omitempty"`
	BitcoinBlockHash       string        `json:"bitcoinblockhash,omitempty"`
}

type JSON struct {
	Raw  string `json:"raw,omitempty"`
	Key  string `json:"key,omitempty"`
	Json string `json:"json,omitempty"`
}

type MerkleNode struct {
	Left  string `json:"left,omitempty"`
	Right string `json:"right,omitempty"`
	Top   string `json:"top,omitempty"`
}
