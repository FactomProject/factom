// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"os"
	"testing"

	"github.com/FactomProject/factom"
)

func TestNewTransaction(t *testing.T) {
	dbpath := os.TempDir() + "/test_wallet-01"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		os.RemoveAll(dbpath)
		t.Error(err)
	}
	
	// create a new transaction
	if err := w1.NewTransaction("test_tx-01"); err != nil {
		t.Error(err)
	}
	if err := w1.NewTransaction("test_tx-02"); err != nil {
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

func TestAddInput (t *testing.T) {
	dbpath := os.TempDir() + "/test_wallet-01"
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		os.RemoveAll(dbpath)
		t.Error(err)
	}
	
	// write a new fct address to the db
	f, err := factom.GetFactoidAddress(zSec)
	if err != nil {
		t.Error(err)
	}
	if err := w1.PutFCTAddress(f); err != nil {
		t.Error(err)
	}
	// Get the address back out of the db
	adr, err := w1.GetFCTAddress(f.String())
	if err != nil {
		t.Error(err)
	}
	
	// create a new transaction
	if err := w1.NewTransaction("tx-01"); err != nil {
		t.Error(err)
	}
	
	if err := w1.AddInput("tx-01", adr.String(), 5); err != nil {
		t.Error(err)
	}
	if len(w1.GetTransactions()) != 1 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}
	t.Log(w1.GetTransactions())

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestComposeTrasnaction(t *testing.T) {
	dbpath := os.TempDir() + "/test_wallet-01"
	f1Sec := "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"
//	f1Sec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	f2Sec := "Fs3GFV6GNV6ar4b8eGcQWpGFbFtkNWKfEPdbywmha8ez5p7XMJyk"
	e1Sec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		os.RemoveAll(dbpath)
		t.Error(err)
	}
	
	// write a new fct address to the db
	if in, err := factom.GetFactoidAddress(f1Sec); err != nil {
		t.Error(err)
	} else {
		if err := w1.PutFCTAddress(in); err != nil {
			t.Error(err)
		}
	}
	
	// Get the address back out of the db
	f1 := factom.NewFactoidAddress()
	if out, err := factom.GetFactoidAddress(f1Sec); err != nil {
		t.Error(err)
	} else {
		if f, err := w1.GetFCTAddress(out.String()); err != nil {
			t.Error(err)
		} else {
			f1 = f
		}
	}
	
	// setup a factoid address for receving
	f2, err := factom.GetFactoidAddress(f2Sec)
	if err != nil {
		t.Error(err)
	}
	
	// setup an ec address for receving
	e1, err := factom.GetECAddress(e1Sec)
	if err != nil {
		t.Error(err)
	}
	
	// create a new transaction
	if err := w1.NewTransaction("tx-01"); err != nil {
		t.Error(err)
	}
	if err := w1.AddInput("tx-01", f1.String(), 5e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddOutput("tx-01", f2.String(), 3e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddECOutput("tx-01", e1.PubString(), 2e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddFee("tx-01", f1.String(), 10000); err != nil {
		t.Error(err)
	}
	if err := w1.SignTransaction("tx-01"); err != nil {
		t.Error(err)
	}
	
	if j, err := w1.ComposeTransaction("tx-01"); err != nil {
		t.Error(err)
	} else {
		t.Log(factom.EncodeJSONString(j))
	}	

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}
