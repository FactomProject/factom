// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type TXInfo struct {
	Name           string `json:"tx-name"`
	TxID           string `json:"txid,omitempty"`
	TotalInputs    uint64 `json:"totalinputs"`
	TotalOutputs   uint64 `json:"totaloutputs"`
	TotalECOutputs uint64 `json:"totalecoutputs"`
	FeesRequired   uint64 `json:"feesrequired,omitempty"`
	RawTransaction string `json:"rawtransaction"`
}

func NewTransaction(name string) error {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("new-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func DeleteTransaction(name string) error {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("delete-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func TransactionHash(name string) (string, error) {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("transaction-hash", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}
	tx := new(TXInfo)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return "", err
	}

	return tx.TxID, nil
}

func ListTransactionsAll() ([]json.RawMessage, error) {
	type transactionList struct {
		Transactions []json.RawMessage `json:"transactions"`
	}

	req := NewJSON2Request("transactions", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(transactionList)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsAddress(addr string) ([]json.RawMessage, error) {
	type transactionList struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	type txReq struct {
		Address string `json:"address"`
	}

	params := txReq{Address: addr}

	req := NewJSON2Request("transactions", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(transactionList)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsID(id string) ([]json.RawMessage, error) {
	type transactionList struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	type txReq struct {
		TxID string `json:"txid"`
	}

	params := txReq{TxID: id}

	req := NewJSON2Request("transactions", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(transactionList)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsRange(start, end int) ([]json.RawMessage, error) {
	type transactionList struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	type txReq struct {
		Range struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"range"`
	}

	params := new(txReq)
	params.Range.Start = start
	params.Range.End = end

	req := NewJSON2Request("transactions", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(transactionList)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsTmp() ([]TXInfo, error) {
	type multiTransactionResponse struct {
		Transactions []TXInfo `json:"transactions"`
	}

	req := NewJSON2Request("tmp-transactions", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	txs := new(multiTransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), txs); err != nil {
		return nil, err
	}
	return txs.Transactions, nil
}

func AddTransactionInput(name, address string, amount uint64) error {
	if AddressStringType(address) != FactoidPub {
		return fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-input", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func AddTransactionOutput(name, address string, amount uint64) error {
	if AddressStringType(address) != FactoidPub {
		return fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-output", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func AddTransactionECOutput(name, address string, amount uint64) error {
	if AddressStringType(address) != ECPub {
		return fmt.Errorf("%s is not an Entry Credit address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-ec-output", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func AddTransactionFee(name, address string) error {
	if AddressStringType(address) != FactoidPub {
		return fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address}
	req := NewJSON2Request("add-fee", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func SubTransactionFee(name, address string) error {
	params := transactionValueRequest{
		Name:    name,
		Address: address}
	req := NewJSON2Request("sub-fee", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func SignTransaction(name string) error {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("sign-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

func ComposeTransaction(name string) ([]byte, error) {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("compose-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.JSONResult(), nil
}

func SendTransaction(name string) (string, error) {
	params := transactionRequest{Name: name}

	wreq := NewJSON2Request("compose-transaction", APICounter(), params)
	wresp, err := walletRequest(wreq)
	if err != nil {
		return "", err
	}
	if wresp.Error != nil {
		return "", wresp.Error
	}

	freq := new(JSON2Request)
	json.Unmarshal(wresp.JSONResult(), freq)
	fresp, err := factomdRequest(freq)
	if err != nil {
		return "", err
	}
	if fresp.Error != nil {
		return "", fresp.Error
	}
	id, err := TransactionHash(name)
	if err != nil {
		return "", err
	}
	if err := DeleteTransaction(name); err != nil {
		return "", err
	}

	return id, nil
}

func SendFactoid(from, to string, amount uint64) (string, error) {
	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return "", err
	}
	name := hex.EncodeToString(n)
	if err := NewTransaction(name); err != nil {
		return "", err
	}
	if err := AddTransactionInput(name, from, amount); err != nil {
		return "", err
	}
	if err := AddTransactionOutput(name, to, amount); err != nil {
		return "", err
	}
	balance, err := GetFactoidBalance(from)
	if err != nil {
		return "", err
	}
	if balance > int64(amount) {
		if err := AddTransactionFee(name, from); err != nil {
			return "", err
		}
	} else {
		if err := SubTransactionFee(name, to); err != nil {
			return "", err
		}
	}
	if err := SignTransaction(name); err != nil {
		return "", err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return "", err
	}

	return r, nil
}

func BuyEC(from, to string, amount uint64) (string, error) {
	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return "", err
	}
	name := hex.EncodeToString(n)
	if err := NewTransaction(name); err != nil {
		return "", err
	}
	if err := AddTransactionInput(name, from, amount); err != nil {
		return "", err
	}
	if err := AddTransactionECOutput(name, to, amount); err != nil {
		return "", err
	}
	if err := AddTransactionFee(name, from); err != nil {
		return "", err
	}
	if err := SignTransaction(name); err != nil {
		return "", err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return "", err
	}

	return r, nil
}

//Purchases the exact amount of ECs
func BuyExactEC(from, to string, amount uint64) (string, error) {
	rate, err := GetRate()
	if err != nil {
		return "", err
	}

	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return "", err
	}
	name := hex.EncodeToString(n)

	if err := NewTransaction(name); err != nil {
		return "", err
	}
	if err := AddTransactionInput(name, from, amount*rate); err != nil {
		return "", err
	}
	if err := AddTransactionECOutput(name, to, amount*rate); err != nil {
		return "", err
	}
	if err := AddTransactionFee(name, from); err != nil {
		return "", err
	}
	if err := SignTransaction(name); err != nil {
		return "", err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return "", err
	}

	return r, nil
}

type TransactionResponse struct {
	ECTranasction      interface{} `json:"ectransaction,omitempty"`
	FactoidTransaction interface{} `json:"factoidtransaction,omitempty"`
	Entry              interface{} `json:"entry,omitempty"`

	//F/EC/E block the transaction is included in
	IncludedInTransactionBlock string `json:"includedintransactionblock"`
	//DirectoryBlock the tranasction is included in
	IncludedInDirectoryBlock string `json:"includedindirectoryblock"`
	//The DBlock height
	IncludedInDirectoryBlockHeight int64 `json:"includedindirectoryblockheight"`
}

func GetTransaction(txID string) (*TransactionResponse, error) {
	params := hashRequest{Hash: txID}
	req := NewJSON2Request("get-transaction", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	txResp := new(TransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), txResp); err != nil {
		return nil, err
	}

	return txResp, nil
}
