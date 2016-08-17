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

type importMnemonicRequest struct {
	Words string `json:"words"`
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

type txdbRequest struct {
	TxID    string `json:"txid,omitempty"`
	Address string `json:"address,omitempty"`
	Range   struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"range,omitempty"`
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
	Name           string `json:"tx-name"`
	TxID           string `json:"txid,omitempty"`
	TotalInputs    uint64 `json:"totalinputs"`
	TotalOutputs   uint64 `json:"totaloutputs"`
	TotalECOutputs uint64 `json:"totalecoutputs"`
	FeesRequired   uint64 `json:"feesrequired,omitempty"`
	RawTransaction string `json:"rawtransaction,omitempty"`
}

type multiTransactionResponse struct {
	Transactions []transactionResponse `json:"transactions"`
}

type propertiesResponse struct {
	WalletVersion string `json:"walletversion"`
}
