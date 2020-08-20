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

// TransAddress represents the imputs and outputs in a Factom Transaction.
type TransAddress struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

// A Transaction from the Factom Network represents a transfer of value between
// Factoid addresses and/or the creation of new Entry Credits.
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

// MarshalJSON converts the Transaction into a JSON object.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	txReq := &struct {
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

	return json.Marshal(txReq)
}

// UnmarshalJSON converts the JSON Transaction back into a Transaction.
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	txResp := new(struct {
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
	})

	if err := json.Unmarshal(data, txResp); err != nil {
		return err
	}

	tx.BlockHeight = txResp.BlockHeight
	tx.FeesPaid = txResp.FeesPaid
	tx.FeesRequired = txResp.FeesRequired
	tx.IsSigned = txResp.IsSigned
	tx.Name = txResp.Name
	tx.Timestamp = time.Unix(txResp.Timestamp, 0)
	tx.TotalECOutputs = txResp.TotalECOutputs
	tx.TotalInputs = txResp.TotalInputs
	tx.TotalOutputs = txResp.TotalOutputs
	tx.Inputs = txResp.Inputs
	tx.Outputs = txResp.Outputs
	tx.ECOutputs = txResp.ECOutputs
	tx.TxID = txResp.TxID

	return nil
}

// NewTransaction creates a new temporary Transaction in the wallet.
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

// DeleteTransaction remove a temporary transacton from the wallet.
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

// ListTransactionsAll lists all the transactions from the wallet database.
func ListTransactionsAll() ([]*Transaction, error) {
	req := NewJSON2Request("transactions", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(struct {
		Transactions []*Transaction `json:"transactions"`
	})
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

// ListTransactionsAddress lists all transaction to and from a given address.
func ListTransactionsAddress(addr string) ([]*Transaction, error) {
	params := &struct {
		Address string `json:"address"`
	}{
		Address: addr,
	}

	req := NewJSON2Request("transactions", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(struct {
		Transactions []*Transaction `json:"transactions"`
	})
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

// ListTransactionsID lists a transaction from the wallet database with a given
// Transaction ID.
func ListTransactionsID(id string) ([]*Transaction, error) {
	params := &struct {
		TxID string `json:"txid"`
	}{
		TxID: id,
	}

	req := NewJSON2Request("transactions", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	list := new(struct {
		Transactions []*Transaction `json:"transactions"`
	})
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

// ListTransactionsRange lists all transacions from the wallet database made
// within a given range of Directory Block heights.
func ListTransactionsRange(start, end int) ([]*Transaction, error) {
	params := new(struct {
		Range struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"range"`
	})
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

	list := new(struct {
		Transactions []*Transaction `json:"transactions"`
	})
	if err := json.Unmarshal(resp.JSONResult(), list); err != nil {
		return nil, err
	}

	return list.Transactions, nil
}

// ListTransactionsTmp lists all of the temporary transaction held in the
// wallet. Temporary transaction are held by the wallet while they are being
// constructed and prepaired to be submitted to the network.
func ListTransactionsTmp() ([]*Transaction, error) {
	req := NewJSON2Request("tmp-transactions", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	txs := new(struct {
		Transactions []*Transaction `json:"transactions"`
	})
	if err := json.Unmarshal(resp.JSONResult(), txs); err != nil {
		return nil, err
	}
	return txs.Transactions, nil
}

// AddTransactionInput adds a factoid input to a temporary transaction in the
// wallet. The imput should come from a Factoid address heald in the wallet
// database.
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
		Amount:  amount,
	}

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

// AddTransactionOutput adds a factoid output to a temporary transaction in
// the wallet.
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
		Amount:  amount,
	}

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

// AddTransactionECOutput adds an Entry Credit output to a temporary transaction
// in the wallet.
func AddTransactionECOutput(name, address string, amount uint64) (*Transaction, error) {
	if AddressStringType(address) != ECPub {
		return nil, fmt.Errorf("%s is not an Entry Credit address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
		Amount:  amount,
	}

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

// AddTransactionFee adds the appropriate factoid fee payment to a transaction
// input of a temporary transaction in the wallet.
func AddTransactionFee(name, address string) (*Transaction, error) {
	if AddressStringType(address) != FactoidPub {
		return nil, fmt.Errorf("%s is not a Factoid address", address)
	}

	params := transactionValueRequest{
		Name:    name,
		Address: address,
	}

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

// SubTransactionFee subtracts the appropriate factoid fee payment from a
// transaction output of a temporary transaction in the wallet.
func SubTransactionFee(name, address string) (*Transaction, error) {
	params := transactionValueRequest{
		Name:    name,
		Address: address,
	}

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

// SignTransaction adds the reqired signatures from the appropriate factoid
// addresses to a temporary transaction in the wallet.
func SignTransaction(name string, force bool) (*Transaction, error) {
	params := transactionRequest{
		Name:  name,
		Force: force,
	}

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

// ComposeTransaction creates a json object from a temporary transaction in the
// wallet that may be sent to the factomd API to submit the transaction to the
// network.
//
// ComposeTransaction may be used by an offline wallet to create an API call
// that can be securely transfered to an online node to enable transactions from
// compleatly offline addresses.
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

// SendTransaction composes a prepaired temoprary transaction from the wallet
// and sends it to the factomd API to be included on the factom network.
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

// SendFactoid creates and sends a transaction to the Factom Network.
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

// BuyEC creates and sends a transaction to the Factom Network that purchases
// Entry Credits.
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

// BuyExactEC creates and sends a transaction to the Factom Network that
// purchases an exact number of Entry Credits.
//
// BuyExactEC calculates the and adds the transaction fees and Entry Credit rate
// so that the exact requested number of Entry Credits are created by the output
// of the transacton.
func BuyExactEC(from, to string, amount uint64, force bool) (*Transaction, error) {
	rate, err := GetECRate()
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

// FactoidSubmit sends a raw transaction to factomd to be included in the
// network. (See ComposeTransaction for more details on how to build the binary
// transaction for the network).
func FactoidSubmit(tx string) (message, txid string, err error) {
	params := &struct {
		Transaction string
	}{
		Transaction: tx,
	}

	req := NewJSON2Request("factoid-submit", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return
	}

	fsr := new(struct {
		Message string `json:"message"`
		TxID    string `json:"txid"`
	})
	if err = json.Unmarshal(resp.JSONResult(), fsr); err != nil {
		return
	}

	return fsr.Message, fsr.TxID, nil
}

// TransactionResponse is the factomd API responce to a request made for a
// transaction.
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

// GetTransaction requests a transaction from the factomd API.
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

// GetTmpTransaction requests a temporary transaction from the wallet.
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
