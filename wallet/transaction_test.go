// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"os"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	dbpath := os.TempDir() + "/test_wallet-01"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		os.RemoveAll(dbpath)
		t.Error(err)
	}
	
	// create a new transaction
	if err := w1.CreateTransaction("test_tx-01"); err != nil {
		t.Error(err)
	}
	if err := w1.CreateTransaction("test_tx-02"); err != nil {
		t.Error(err)
	}
	
	if len(w1.GetTransactions()) != 2 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}
	
	if err := w1.DeleteTransaction("test_tx-02"); err != nil {
		t.Error(err)
	}

	if len(w1.GetTransactions()) != 1 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}
