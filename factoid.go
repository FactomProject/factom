// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

// FactoidSubmit sends a transaction to factomd to be included in the network.
// (See ComposeTransaction for more details on how to build the binary
// transaction for the network).
func FactoidSubmit(tx string) (message, txid string, err error) {
	type txreq struct {
		Transaction string
	}

	params := txreq{Transaction: tx}
	req := NewJSON2Request("factoid-submit", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return
	}

	fsr := new(struct {
		Message string `json:"message"`
		TxID    string `json:"txid"`
	})
	if err = json.Unmarshal(resp.JSONResult(), fsr); err != nil {
		return
	}

	return fsr.Message, fsr.TxID, nil
}
