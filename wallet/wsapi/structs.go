// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

import (
	"github.com/FactomProject/factom"
)

type TLSConfig struct {
	TLSEnable   bool
	TLSKeyFile  string
	TLSCertFile string
}

// requests

type addressRequest struct {
	Address string `json:"address"`
}

type importRequest struct {
	Addresses []struct {
		Secret string `json:"secret"`
	} `json:addresses`
}

type importKoinifyRequest struct {
	Words string `json:"words"`
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

type txdbRequest struct {
	TxID    string `json:"txid,omitempty"`
	Address string `json:"address,omitempty"`
	Range   struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"range,omitempty"`
}

type entryRequest struct {
	Entry factom.Entry `json:"entry"`
	ECPub string       `json:"ecpub"`
	Force bool         `json:"force"`
}

type chainRequest struct {
	Chain factom.Chain `json:"chain"`
	ECPub string       `json:"ecpub"`
	Force bool         `json:"force"`
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

type multiTransactionResponse struct {
	Transactions []*factom.Transaction `json:"transactions"`
}

type propertiesResponse struct {
	WalletVersion    string `json:"walletversion"`
	WalletApiVersion string `json:"walletapiversion"`
}

type simpleResponse struct {
	Success bool `json:"success"`
}

type entryResponse struct {
	Commit *factom.JSON2Request `json:"commit"`
	Reveal *factom.JSON2Request `json:"reveal"`
}

type heightResponse struct {
	Height int64 `json:"height"`
}
