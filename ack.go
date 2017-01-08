// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

type FactoidTxStatus struct {
	TxID string `json:"txid"`
	GeneralTransactionData
}

func (f *FactoidTxStatus) String() string {
	var s string
	s += fmt.Sprintln("TxID:", f.TxID)
	s += fmt.Sprintln("Status:", f.Status)
	s += fmt.Sprintln("Date:", f.TransactionDateString)

	return s
}

type EntryStatus struct {
	CommitTxID string `json:"committxid"`
	EntryHash  string `json:"entryhash"`

	CommitData GeneralTransactionData `json:"commitdata"`
	EntryData  GeneralTransactionData `json:"entrydata"`

	ReserveTransactions          []ReserveInfo `json:"reserveinfo,omitempty"`
	ConflictingRevealEntryHashes []string      `json:"conflictingrevealentryhashes,omitempty"`
}

func (e *EntryStatus) String() string {
	var s string
	s += fmt.Sprintln("TxID:", e.CommitTxID)
	s += fmt.Sprintln("Status:", e.CommitData.Status)
	s += fmt.Sprintln("Date:", e.CommitData.TransactionDateString)

	return s
}

type ReserveInfo struct {
	TxID    string `json:"txid"`
	Timeout int64  `json:"timeout"` //Unix time
}

type GeneralTransactionData struct {
	// TransactionDate in Unix time
	TransactionDate int64 `json:"transactiondate,omitempty"`
	//TransactionDateString ISO8601 time
	TransactionDateString string `json:"transactiondatestring,omitempty"`
	//Unix time
	BlockDate int64 `json:"blockdate,omitempty"`
	//ISO8601 time
	BlockDateString string `json:"blockdatestring,omitempty"`

	Malleated *Malleated `json:"malleated,omitempty"`
	Status    string     `json:"status"`
}

type Malleated struct {
	MalleatedTxIDs []string `json:"malleatedtxids"`
}

func FactoidACK(txID, fullTransaction string) (*FactoidTxStatus, error) {
	params := ackRequest{TxID: txID, FullTransaction: fullTransaction}
	req := NewJSON2Request("factoid-ack", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(FactoidTxStatus)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}

func EntryACK(txID, fullTransaction string) (*EntryStatus, error) {
	params := ackRequest{TxID: txID, FullTransaction: fullTransaction}
	req := NewJSON2Request("entry-ack", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(EntryStatus)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}
