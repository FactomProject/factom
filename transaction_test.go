// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"encoding/json"
	"time"

	. "github.com/FactomProject/factom"

	"testing"
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
