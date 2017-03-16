// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	. "github.com/FactomProject/factom"
	"testing"

	"encoding/json"
	"time"
)

func TestJSONTransactions(t *testing.T) {
	tx1 := mkdummytx()
	t.Log("Transaction:", tx1)
	p, err := json.Marshal(tx1)
	if err != nil {
		t.Error(err)
	}
	t.Log("JSON transaction:", string(p))

	tx2 := new(Transaction)
	if err := json.Unmarshal(p, tx2); err != nil {
		t.Error(err)
	}
	t.Log("Unmarshaled:", tx2)
}

func TestTransactions(t *testing.T) {
	// start the test wallet
	done, err := StartTestWallet()
	if err != nil {
		t.Error(err)
	}
	defer func() { done <- 1 }()

	// make sure the wallet is empty
	if txs, err := ListTransactionsTmp(); err != nil {
		t.Error(err)
	} else if len(txs) > 0 {
		t.Error("Unexpected transactions returned from the wallet:", txs)
	}

	// create a new transaction
	tx1, err := NewTransaction("tx1")
	if err != nil {
		t.Error(err)
	}
	if tx1 == nil {
		t.Error("No transaction was returned")
	}

	if tx, err := GetTmpTransaction("tx1"); err != nil {
		t.Error(err)
	} else if tx == nil {
		t.Error("Temporary transaction was not saved in the wallet")
	}

	// delete a transaction
	if err := DeleteTransaction("tx1"); err != nil {
		t.Error(err)
	}

	if txs, err := ListTransactionsTmp(); err != nil {
		t.Error(err)
	} else if len(txs) > 0 {
		t.Error("Unexpected transactions returned from the wallet:", txs)
	}
}

// helper functions for testing

func mkdummytx() *Transaction {
	tx := &Transaction{
		BlockHeight: 42,
		Name:        "dummy",
		Timestamp: func() time.Time {
			t, _ := time.Parse("2006-Jan-02 15:04", "1988-Jan-02 10:00")
			return t
		}(),
		TotalInputs:    13,
		TotalOutputs:   12,
		TotalECOutputs: 1,
	}
	return tx
}
