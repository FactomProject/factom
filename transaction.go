// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type TransAddress struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type Transaction struct {
	BlockHeight    uint32          `json:"blockheight,omitempty"`
	FeesPaid       uint64          `json:"feespaid,omitempty"`
	FeesRequired   uint64          `json:"feesrequired,omitempty"`
	IsSigned       bool            `json:"signed"`
	Name           string          `json:"name,omitempty"`
	Timestamp      time.Time       `json:"timestamp"`
	TotalECOutputs uint64          `json:"totalecoutputs"`
	TotalInputs    uint64          `json:"totalinputs"`
	TotalOutputs   uint64          `json:"totaloutputs"`
	Inputs         []*TransAddress `json:"inputs"`
	Outputs        []*TransAddress `json:"outputs"`
	ECOutputs      []*TransAddress `json:"ecoutputs"`
	TxID           string          `json:"txid,omitempty"`
}

// String prints the formatted data of a transaction.
func (tx *Transaction) String() (s string) {
	if tx.Name != "" {
		s += fmt.Sprintln("Name:", tx.Name)
	}
	if tx.IsSigned {
		s += fmt.Sprintln("TxID:", tx.TxID)
	}
	s += fmt.Sprintln("Timestamp:", tx.Timestamp)
	if tx.BlockHeight != 0 {
		s += fmt.Sprintln("BlockHeight:", tx.BlockHeight)
	}
	s += fmt.Sprintln("TotalInputs:", FactoshiToFactoid(tx.TotalInputs))
	s += fmt.Sprintln("TotalOutputs:", FactoshiToFactoid(tx.TotalOutputs))
	s += fmt.Sprintln("TotalECOutputs:", FactoshiToFactoid(tx.TotalECOutputs))
	for _, in := range tx.Inputs {
		s += fmt.Sprintln(
			"Input:",
			in.Address,
			FactoshiToFactoid(in.Amount),
		)
	}
	for _, out := range tx.Outputs {
		s += fmt.Sprintln(
			"Output:",
			out.Address,
			FactoshiToFactoid(out.Amount),
		)
	}
	for _, ec := range tx.ECOutputs {
		s += fmt.Sprintln(
			"ECOutput:",
			ec.Address,
			FactoshiToFactoid(ec.Amount),
		)
	}
	s += fmt.Sprintln("FeesPaid:", FactoshiToFactoid(tx.FeesPaid))
	if tx.FeesRequired != 0 {
		s += fmt.Sprintln("FeesRequired:", FactoshiToFactoid(tx.FeesRequired))
	}
	s += fmt.Sprintln("Signed:", tx.IsSigned)

	return s
}

// MarshalJSON converts the Transaction into a JSON object
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	tmp := &struct {
		BlockHeight    uint32          `json:"blockheight,omitempty"`
		FeesPaid       uint64          `json:"feespaid,omitempty"`
		FeesRequired   uint64          `json:"feesrequired,omitempty"`
		IsSigned       bool            `json:"signed"`
		Name           string          `json:"name,omitempty"`
		Timestamp      int64           `json:"timestamp"`
		TotalECOutputs uint64          `json:"totalecoutputs"`
		TotalInputs    uint64          `json:"totalinputs"`
		TotalOutputs   uint64          `json:"totaloutputs"`
		Inputs         []*TransAddress `json:"inputs"`
		Outputs        []*TransAddress `json:"outputs"`
		ECOutputs      []*TransAddress `json:"ecoutputs"`
		TxID           string          `json:"txid,omitempty"`
	}{
		BlockHeight:    tx.BlockHeight,
		FeesPaid:       tx.FeesPaid,
		FeesRequired:   tx.FeesRequired,
		IsSigned:       tx.IsSigned,
		Name:           tx.Name,
		Timestamp:      tx.Timestamp.Unix(),
		TotalECOutputs: tx.TotalECOutputs,
		TotalInputs:    tx.TotalInputs,
		TotalOutputs:   tx.TotalOutputs,
		Inputs:         tx.Inputs,
		Outputs:        tx.Outputs,
		ECOutputs:      tx.ECOutputs,
		TxID:           tx.TxID,
	}

	return json.Marshal(tmp)
}

// UnmarshalJSON converts the JSON Transaction back into a Transaction
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	type jsontx struct {
		BlockHeight    uint32          `json:"blockheight,omitempty"`
		FeesPaid       uint64          `json:"feespaid,omitempty"`
		FeesRequired   uint64          `json:"feesrequired,omitempty"`
		IsSigned       bool            `json:"signed"`
		Name           string          `json:"name,omitempty"`
		Timestamp      int64           `json:"timestamp"`
		TotalECOutputs uint64          `json:"totalecoutputs"`
		TotalInputs    uint64          `json:"totalinputs"`
		TotalOutputs   uint64          `json:"totaloutputs"`
		Inputs         []*TransAddress `json:"inputs"`
		Outputs        []*TransAddress `json:"outputs"`
		ECOutputs      []*TransAddress `json:"ecoutputs"`
		TxID           string          `json:"txid,omitempty"`
	}
	tmp := new(jsontx)

	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}

	tx.BlockHeight = tmp.BlockHeight
	tx.FeesPaid = tmp.FeesPaid
	tx.FeesRequired = tmp.FeesRequired
	tx.IsSigned = tmp.IsSigned
	tx.Name = tmp.Name
	tx.Timestamp = time.Unix(tmp.Timestamp, 0)
	tx.TotalECOutputs = tmp.TotalECOutputs
	tx.TotalInputs = tmp.TotalInputs
	tx.TotalOutputs = tmp.TotalOutputs
	tx.Inputs = tmp.Inputs
	tx.Outputs = tmp.Outputs
	tx.ECOutputs = tmp.ECOutputs
	tx.TxID = tmp.TxID

	return nil
}

// NewTransaction creates a new temporary Transaction in the wallet
func NewTransaction(name string) (*Transaction, error) {
	params := transactionRequest{Name: name}
	req := NewJSON2Request("new-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
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

func ListTransactionsAll() ([]*Transaction, error) {
	type multiTransactionResponse struct {
		Transactions []*Transaction `json:"transactions"`
	}

	req := NewJSON2Request("transactions", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(multiTransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsAddress(addr string) ([]*Transaction, error) {
	type multiTransactionResponse struct {
		Transactions []*Transaction `json:"transactions"`
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

	list := new(multiTransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsID(id string) ([]*Transaction, error) {
	type multiTransactionResponse struct {
		Transactions []*Transaction `json:"transactions"`
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

	list := new(multiTransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsRange(start, end int) ([]*Transaction, error) {
	type multiTransactionResponse struct {
		Transactions []*Transaction `json:"transactions"`
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

	list := new(multiTransactionResponse)
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

func ListTransactionsTmp() ([]*Transaction, error) {
	type multiTransactionResponse struct {
		Transactions []*Transaction `json:"transactions"`
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

func AddTransactionInput(
	name,
	address string,
	amount uint64,
) (*Transaction, error) {
	if AddressStringType(address) != FactoidPub {
		return nil, fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-input", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func AddTransactionOutput(
	name,
	address string,
	amount uint64,
) (*Transaction, error) {
	if AddressStringType(address) != FactoidPub {
		return nil, fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-output", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func AddTransactionECOutput(
	name,
	address string,
	amount uint64,
) (*Transaction, error) {
	if AddressStringType(address) != ECPub {
		return nil, fmt.Errorf("%s is not an Entry Credit address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount}
	req := NewJSON2Request("add-ec-output", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func AddTransactionFee(name, address string) (*Transaction, error) {
	if AddressStringType(address) != FactoidPub {
		return nil, fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address}
	req := NewJSON2Request("add-fee", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func SubTransactionFee(name, address string) (*Transaction, error) {
	params := transactionValueRequest{
		Name:    name,
		Address: address}
	req := NewJSON2Request("sub-fee", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func SignTransaction(name string, force bool) (*Transaction, error) {
	params := transactionRequest{Name: name}
	params.Force = force
	req := NewJSON2Request("sign-transaction", APICounter(), params)

	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	tx := new(Transaction)
	if err := json.Unmarshal(resp.JSONResult(), tx); err != nil {
		return nil, err
	}
	return tx, nil
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

func SendTransaction(name string) (*Transaction, error) {
	params := transactionRequest{Name: name}

	tx, err := GetTmpTransaction(name)
	if err != nil {
		return nil, err
	}
	if !tx.IsSigned {
		return nil, fmt.Errorf("Cannot send unsigned transaction")
	}

	wreq := NewJSON2Request("compose-transaction", APICounter(), params)
	wresp, err := walletRequest(wreq)
	if err != nil {
		return nil, err
	}
	if wresp.Error != nil {
		return nil, wresp.Error
	}

	freq := new(JSON2Request)
	json.Unmarshal(wresp.JSONResult(), freq)
	fresp, err := factomdRequest(freq)
	if err != nil {
		return nil, err
	}
	if fresp.Error != nil {
		return nil, fresp.Error
	}
	if err := DeleteTransaction(name); err != nil {
		return nil, err
	}

	return tx, nil
}

func SendFactoid(from, to string, amount uint64, force bool) (*Transaction, error) {
	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return nil, err
	}
	name := hex.EncodeToString(n)
	if _, err := NewTransaction(name); err != nil {
		return nil, err
	}
	if _, err := AddTransactionInput(name, from, amount); err != nil {
		return nil, err
	}
	if _, err := AddTransactionOutput(name, to, amount); err != nil {
		return nil, err
	}
	balance, err := GetFactoidBalance(from)
	if err != nil {
		return nil, err
	}
	if balance > int64(amount) {
		if _, err := AddTransactionFee(name, from); err != nil {
			return nil, err
		}
	} else {
		if _, err := SubTransactionFee(name, to); err != nil {
			return nil, err
		}
	}
	if _, err := SignTransaction(name, force); err != nil {
		return nil, err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func BuyEC(from, to string, amount uint64, force bool) (*Transaction, error) {
	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return nil, err
	}
	name := hex.EncodeToString(n)
	if _, err := NewTransaction(name); err != nil {
		return nil, err
	}
	if _, err := AddTransactionInput(name, from, amount); err != nil {
		return nil, err
	}
	if _, err := AddTransactionECOutput(name, to, amount); err != nil {
		return nil, err
	}
	if _, err := AddTransactionFee(name, from); err != nil {
		return nil, err
	}
	if _, err := SignTransaction(name, force); err != nil {
		return nil, err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return nil, err
	}

	return r, nil
}

//Purchases the exact amount of ECs
func BuyExactEC(from, to string, amount uint64, force bool) (*Transaction, error) {
	rate, err := GetRate()
	if err != nil {
		return nil, err
	}

	n := make([]byte, 16)
	if _, err := rand.Read(n); err != nil {
		return nil, err
	}
	name := hex.EncodeToString(n)

	if _, err := NewTransaction(name); err != nil {
		return nil, err
	}
	if _, err := AddTransactionInput(name, from, amount*rate); err != nil {
		return nil, err
	}
	if _, err := AddTransactionECOutput(name, to, amount*rate); err != nil {
		return nil, err
	}
	if _, err := AddTransactionFee(name, from); err != nil {
		return nil, err
	}
	if _, err := SignTransaction(name, force); err != nil {
		return nil, err
	}
	r, err := SendTransaction(name)
	if err != nil {
		return nil, err
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
	req := NewJSON2Request("transaction", APICounter(), params)
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

// GetTmpTransaction gets a temporary transaction from the wallet
func GetTmpTransaction(name string) (*Transaction, error) {
	txs, err := ListTransactionsTmp()
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		if tx.Name == name {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("Transaction not found")
}
