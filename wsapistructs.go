// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import "fmt"

// requests

type heightRequest struct {
	Height int64 `json:"height"`
}

type ackRequest struct {
	Hash            string `json:"hash,omitempty"`
	ChainID         string `json:"chainid,omitempty"`
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

type HeightsResponse struct {
	DirectoryBlockHeight int64 `json:"directoryblockheight"`
	LeaderHeight         int64 `json:"leaderheight"`
	EntryBlockHeight     int64 `json:"entryblockheight"`
	EntryHeight          int64 `json:"entryheight"`
}

func (d *HeightsResponse) String() string {
	var s string

	s += fmt.Sprintln("DirectoryBlockHeight:", d.DirectoryBlockHeight)
	s += fmt.Sprintln("LeaderHeight:", d.LeaderHeight)
	s += fmt.Sprintln("EntryBlockHeight:", d.EntryBlockHeight)
	s += fmt.Sprintln("EntryHeight:", d.EntryHeight)

	return s
}

type importRequest struct {
	Addresses []secretRequest `json:"addresses"`
}

type importKoinifyRequest struct {
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
