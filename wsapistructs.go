// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import ()

// requests

type ackRequest struct {
	TxID            string `json:"txid,omitempty"`
	FullTransaction string `json:"fulltransaction,omitempty"`
}

type addressRequest struct {
	Address string `json:"address"`
}

type chainIDRequest struct {
	ChainID string `json:"chainid"`
}

type entryRequest struct {
	Entry string `json:"entry"`
}

type hashRequest struct {
	Hash string `json:"hash"`
}

type importRequest struct {
	Addresses []secretRequest `json:"addresses"`
}

type importMnemonicRequest struct {
	Words string `json:"words"`
}

type keyMRRequest struct {
	KeyMR string `json:"keymr"`
}

type messageRequest struct {
	Message string `json:"message"`
}

type secretRequest struct {
	Secret string `json:"secret"`
}

type transactionRequest struct {
	Name string `json:"tx-name"`
}

type transactionValueRequest struct {
	Name    string `json:"tx-name"`
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type transactionAddressRequest struct {
	Name    string `json:"tx-name"`
	Address string `json:"address"`
}
