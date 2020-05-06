// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

// requests

type heightRequest struct {
	Height int64 `json:"height"`
	NoRaw  bool  `json:"noraw,omitempty"`
}

type replayRequest struct {
	StartHeight int64 `json:"startheight"`
	EndHeight   int64 `json:"endheight,omitempty"`
}

type replayResponse struct {
	Message string `json:"message"`
	Start   int64  `json:"startheight"`
	End     int64  `json:"endheight"`
}

type ackRequest struct {
	Hash            string `json:"hash,omitempty"`
	ChainID         string `json:"chainid,omitempty"`
	FullTransaction string `json:"fulltransaction,omitempty"`
}

type addressRequest struct {
	Address string `json:"address"`
}

type passphraseRequest struct {
	Password string `json:"passphrase"`
	Timeout  int64  `json:"timeout"`
}

type unlockResponse struct {
	Success       bool  `json:"success"`
	UnlockedUntil int64 `json:"unlockeduntil"`
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

type keyMRRequest struct {
	KeyMR string `json:"keymr"`
	NoRaw bool   `json:"noraw,omitempty"`
}

type messageRequest struct {
	Message string `json:"message"`
}

type secretRequest struct {
	Secret string `json:"secret"`
}

type transactionRequest struct {
	Name  string `json:"tx-name"`
	Force bool   `json:"force"`
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
