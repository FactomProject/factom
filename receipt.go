// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

func GetReceipt(hash string) (*Receipt, error) {
	type receiptResponse struct {
		Receipt *Receipt `json:"receipt"`
	}

	params := hashRequest{Hash: hash}
	req := NewJSON2Request("receipt", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	rec := new(receiptResponse)
	if err := json.Unmarshal(resp.JSONResult(), rec); err != nil {
		return nil, err
	}

	return rec.Receipt, nil
}

type Receipt struct {
	Entry struct {
		Raw       string `json:"raw,omitempty"`
		EntryHash string `json:"entryhash,omitempty"`
		Json      string `json:"json,omitempty"`
	} `json:"entry,omitempty"`
	MerkleBranch []struct {
		Left  string `json:"left,omitempty"`
		Right string `json:"right,omitempty"`
		Top   string `json:"top,omitempty"`
	} `json:"merklebranch,omitempty"`
	EntryBlockKeyMR        string `json:"entryblockkeymr,omitempty"`
	DirectoryBlockKeyMR    string `json:"directoryblockkeymr,omitempty"`
	BitcoinTransactionHash string `json:"bitcointransactionhash,omitempty"`
	BitcoinBlockHash       string `json:"bitcoinblockhash,omitempty"`
}
