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
	RawTransaction string `json:"rawtransaction"`
}

func NewTransaction(name string) error {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("new-transaction", apiCounter(), params)

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
	req := NewJSON2Request("delete-transaction", apiCounter(), params)

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
	req := NewJSON2Request("transaction-hash", apiCounter(), params)

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

func ListTransactions() ([]TXInfo, error) {
	type multiTransactionResponse struct {
		Transactions []TXInfo `json:"transactions"`
	}

	req := NewJSON2Request("transactions", apiCounter(), nil)
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
	req := NewJSON2Request("add-input", apiCounter(), params)

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
	req := NewJSON2Request("add-output", apiCounter(), params)

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
	req := NewJSON2Request("add-ec-output", apiCounter(), params)

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
	req := NewJSON2Request("add-fee", apiCounter(), params)

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
	req := NewJSON2Request("sub-fee", apiCounter(), params)

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
	req := NewJSON2Request("sign-transaction", apiCounter(), params)

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
	req := NewJSON2Request("compose-transaction", apiCounter(), params)

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

	wreq := NewJSON2Request("compose-transaction", apiCounter(), params)
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

func BuyEC(from, to string, ammount uint64) (string, error) {
	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return "", err
	}
	name := hex.EncodeToString(n)
	if err := NewTransaction(name); err != nil {
		return "", err
	}
	if err := AddTransactionInput(name, from, ammount); err != nil {
		return "", err
	}
	if err := AddTransactionECOutput(name, to, ammount); err != nil {
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
