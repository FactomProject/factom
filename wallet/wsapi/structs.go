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

type passphraseRequest struct {
	Password string `json:"passphrase"`
	Timeout  int64  `json:"timeout"`
}

type addressRequest struct {
	Address string `json:"address"`
}

type addressesRequest struct {
	Addresses []string `json:"addresses"`
}

type importRequest struct {
	Addresses []struct {
		Secret string `json:"secret"`
	} `json:"addresses"`
}

type importKoinifyRequest struct {
	Words string `json:"words"`
}

type transactionRequest struct {
	Name  string `json:"tx-name"`
	Force bool   `json:"force"`
}

type signDataRequest struct {
	Signer string `json:"signer"`
	Data   []byte `json:"data"`
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

type identityKeyRequest struct {
	Public string `json:"public"`
}

type importIdentityKeysRequest struct {
	Keys []struct {
		Secret string `json:"secret"`
	} `json:keys`
}

type activeIdentityKeysRequest struct {
	ChainID string `json:"chainid"`
	Height  *int64 `json:"height"`
}

type identityChainRequest struct {
	Name    []string `json:"name"`
	PubKeys []string `json:"pubkeys"`
	ECPub   string   `json:"ecpub"`
	Force   bool     `json:"force"`
}

type identityKeyReplacementRequest struct {
	ChainID   string `json:"chainid"`
	OldKey    string `json:"oldkey"`
	NewKey    string `json:"newkey"`
	SignerKey string `json:"signerkey"`
	ECPub     string `json:"ecpub"`
	Force     bool   `json:"force"`
}

type identityAttributeRequest struct {
	ReceiverChainID    string                     `json:"receiver-chainid"`
	DestinationChainID string                     `json:"destination-chainid"`
	Attributes         []factom.IdentityAttribute `json:"attributes"`
	SignerKey          string                     `json:"signerkey"`
	SignerChainID      string                     `json:"signer-chainid"`
	ECPub              string                     `json:"ecpub"`
	Force              bool                       `json:"force"`
}

type identityAttributeEndorsementRequest struct {
	DestinationChainID string `json:"destination-chainid"`
	EntryHash          string `json:"entry-hash"`
	SignerKey          string `json:"signerkey"`
	SignerChainID      string `json:"signer-chainid"`
	ECPub              string `json:"ecpub"`
	Force              bool   `json:"force"`
}

// responses

type addressResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

type multiAddressResponse struct {
	Addresses []*addressResponse `json:"addresses"`
}

type balanceResponse struct {
	CurrentHeight   uint32        `json:"current-height"`
	LastSavedHeight uint          `json:"last-saved-height"`
	Balances        []interface{} `json:"balances"`
}

type multiBalanceResponse struct {
	FactoidAccountBalances struct {
		Ack   int64 `json:"ack"`
		Saved int64 `json:"saved"`
	} `json:"fctaccountbalances"`
	EntryCreditAccountBalances struct {
		Ack   int64 `json:"ack"`
		Saved int64 `json:"saved"`
	} `json:"ecaccountbalances"`
}

type walletBackupResponse struct {
	Seed         string                 `json:"wallet-seed"`
	Addresses    []*addressResponse     `json:"addresses"`
	IdentityKeys []*identityKeyResponse `json:"identity-keys"`
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

type unlockResponse struct {
	Success       bool  `json:"success"`
	UnlockedUntil int64 `json:"unlockeduntil"`
}

type entryResponse struct {
	Commit *factom.JSON2Request `json:"commit"`
	Reveal *factom.JSON2Request `json:"reveal"`
}

type heightResponse struct {
	Height int64 `json:"height"`
}

type identityKeyResponse struct {
	Public string `json:"public"`
	Secret string `json:"secret,omitempty"`
}

type multiIdentityKeyResponse struct {
	Keys []*identityKeyResponse `json:"keys"`
}

type activeIdentityKeysResponse struct {
	ChainID string   `json:"chainid"`
	Height  int64    `json:"height"`
	Keys    []string `json:"keys"`
}

type signDataResponse struct {
	PubKey    []byte `json:"pubkey"`
	Signature []byte `json:"signature"`
}

// Helper structs

type UnmarBody struct {
	Jsonrpc string          `json:"jsonrps"`
	Id      int             `json:"id"`
	Result  balanceResponse `json:"result"`
}
