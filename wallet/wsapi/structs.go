// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

// requests

type addressRequest struct {
	Address string `json:"address"`
}

type importRequest struct {
	Addresses []struct {
		Secret string `json:"secret"`
	} `json:addresses`
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

// responses

type addressResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type multiAddressResponse struct {
	Addresses []*addressResponse `json:"addresses"`
}

type walletBackupResponse struct {
	Seed      string             `json:"wallet-seed"`
	Addresses []*addressResponse `json:"addresses"`
}

type transactionResponse struct {
	Name        string `json:"tx-name"`
	Transaction string `json:"transaction"`
}

type transactionsResponse struct {
	Transactions []transactionResponse `json:"transactions"`
}
