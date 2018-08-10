// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// FactoidTxStatus embeds GeneralTransactionData with a transaction ID.
type FactoidTxStatus struct {
	TxID string `json:"txid"`
	GeneralTransactionData
}

// String satisfies the Stringer interface for FactoidTxStatus and returns a
// formatted string containing TxID, Status, and TransactionDateString.
func (f *FactoidTxStatus) String() string {
	var s string
	s += fmt.Sprintln("TxID:", f.TxID)
	s += fmt.Sprintln("Status:", f.Status)
	s += fmt.Sprintln("Date:", f.TransactionDateString)

	return s
}

// EntryStatus represents metadata about an Entry that has been committed and
// revealed.
type EntryStatus struct {
	CommitTxID string `json:"committxid"`
	EntryHash  string `json:"entryhash"`

	CommitData GeneralTransactionData `json:"commitdata"`
	EntryData  GeneralTransactionData `json:"entrydata"`

	ReserveTransactions          []ReserveInfo `json:"reserveinfo,omitempty"`
	ConflictingRevealEntryHashes []string      `json:"conflictingrevealentryhashes,omitempty"`
}

// String satisfies the Stringer interface for EntryStatus and returns a
// formatted string containing EntryHash, EntryData.Status,
// EntryData.TransactionDateString, CommitTxID, CommitData.Status,
// CommitData.TransactionDateString.
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

// ReserveInfo associates timeout information with a transaction ID.
type ReserveInfo struct {
	TxID    string `json:"txid"`
	Timeout int64  `json:"timeout"` //Unix time
}

// GeneralTransactionData contains metadata about a transaction.
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

// Malleated is a struct that wraps a slice of malleated transaction IDs.
type Malleated struct {
	MalleatedTxIDs []string `json:"malleatedtxids"`
}

// EntryCommitACK makes an "ack" RPC request using the txID of the commit to
// search for the entry OR chain commit.
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

// FactoidACK makes an "ack" RPC request to check the status of a Factoid
// transaction.
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

// EntryRevealACK makes an "ack" RPC request using the entryhash to search for
// the entry or chain reveal.
func EntryRevealACK(entryhash, fullTransaction, chainID string) (*EntryStatus, error) {
	params := ackRequest{Hash: entryhash, ChainID: chainID, FullTransaction: fullTransaction}
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
// Use either EntryCommitAck or EntryRevealAck depending on the type of hash
// you are sending.
func EntryACK(entryhash, fullTransaction string) (*EntryStatus, error) {
	return EntryRevealACK(entryhash, fullTransaction, "0000000000000000000000000000000000000000000000000000000000000000")
}
