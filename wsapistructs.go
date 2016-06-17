// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import ()

// requests

type nameRequest struct {
	Name string `json:"name"`
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

type secretRequest struct {
	Secret string `json:"secret"`
}

type importRequest struct {
	Addresses []secretRequest `json:addresses`
}

type keyMRRequest struct {
	KeyMR string `json:"keymr"`
}

type keyRequest struct {
	Key string `json:"key"`
}

type messageRequest struct {
	Message string `json:"message"`
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
