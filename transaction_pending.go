package factom

import (
	"encoding/json"
)

// PendingTransaction is a single transaction returned by the pending-transaction API call
type PendingTransaction struct {
	TxID      string `json:"transactionid"`
	Status    string `json:"status"`
	Fees      uint64 `json:"fees"`
	Inputs    []PendingAddress
	Outputs   []PendingAddress
	ECOutputs []PendingAddress
}

// PendingAddress is the input or recipient of a transaction
type PendingAddress struct {
	Amount  uint64 `json:"amount"`
	RCDHash string `json:"address"`
	Address string `json:"useraddress"`
}

// GetPendingTransactions requests a list of transactions that have been
// submitted to the Factom Network, but have not yet been included in a Factoid
// Block.
func GetPendingTransactions() ([]PendingTransaction, error) {
	req := NewJSON2Request("pending-transactions", APICounter(), nil)
	resp, err := factomdRequest(req)

	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, err
	}

	transList := resp.JSONResult()

	var res []PendingTransaction
	if err := json.Unmarshal(transList, &res); err != nil {
		return nil, err
	}

	return res, nil
}
