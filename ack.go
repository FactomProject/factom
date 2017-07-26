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
	if e.EntryHash != "" {
		s += fmt.Sprintln("EntryHash:", e.EntryHash)
		s += fmt.Sprintln("Status:", e.EntryData.Status)
		s += fmt.Sprintln("Date:", e.EntryData.TransactionDateString)
	}
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

// EntryCommitACK takes the txid of the commit and searches for the entry/chain commit
func EntryCommitACK(txID, fullTransaction string) (*EntryStatus, error) {
	params := ackRequest{Hash: txID, ChainID: "c", FullTransaction: fullTransaction}
	req := NewJSON2Request("ack", APICounter(), params)
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

func FactoidACK(txID, fullTransaction string) (*FactoidTxStatus, error) {
	params := ackRequest{Hash: txID, ChainID: "f", FullTransaction: fullTransaction}
	req := NewJSON2Request("ack", APICounter(), params)
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

// EntryRevealACK will take the entryhash and search for the entry and the commit
func EntryRevealACK(entryhash, fullTransaction, chainiID string) (*EntryStatus, error) {
	params := ackRequest{Hash: entryhash, ChainID: chainiID, FullTransaction: fullTransaction}
	req := NewJSON2Request("ack", APICounter(), params)
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

// EntryACK is a deprecated call and SHOULD NOT BE USED.
// Use either EntryCommitAck or EntryRevealAck depending on the
// type of hash you are sending.
func EntryACK(entryhash, fullTransaction string) (*EntryStatus, error) {
	return EntryRevealACK(entryhash, fullTransaction, "0000000000000000000000000000000000000000000000000000000000000000")
}
