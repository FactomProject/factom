// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	"os"
	"testing"

	"github.com/FactomProject/factom"
	. "github.com/FactomProject/factom/wallet"
)

func TestNewWallet(t *testing.T) {
	// create a new database
	w1, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// check that the seed got written
	seed, err := w1.GetDBSeed()
	if err != nil {
		t.Error(err)
	}
	if len(seed.MnemonicSeed) == 0 {
		t.Errorf("stored db seed is empty")
	}

	if err := w1.Close(); err != nil {
		t.Error(err)
	}
}

func TestOpenWallet(t *testing.T) {
	dbpath := os.TempDir() + "/test_wallet-01"

	// create a new database
	w1, err := NewOrOpenLevelDBWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	w1.Close()

	// make sure we can open the db
	w2, err := NewOrOpenLevelDBWallet(dbpath)
	if err != nil {
		t.Error(err)
	}

	// check that the seed is there
	seed, err := w1.GetDBSeed()
	if len(seed.MnemonicSeed) == 0 {
		t.Errorf("stored db seed is empty")
	}

	if err := w2.Close(); err != nil {
		t.Error(err)
	}

	// remove the testing db
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestPutECAddress(t *testing.T) {
	zSec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"

	// create a new database
	w, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// write a new ec address to the db
	e, err := factom.GetECAddress(zSec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertECAddress(e); err != nil {
		t.Error(err)
	}

	// Check that the address was written into the db
	if _, err := w.GetECAddress(e.PubString()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

func TestPutFCTAddress(t *testing.T) {
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"

	// create a new database
	w, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// write a new fct address to the db
	f, err := factom.GetFactoidAddress(zSec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertFCTAddress(f); err != nil {
		t.Error(err)
	}

	// Check that the address was written into the db
	if _, err := w.GetFCTAddress(f.String()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

func TestGenerateECAddress(t *testing.T) {
	// create a new database
	w, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// Generate a new ec address
	e, err := w.GenerateECAddress()
	if err != nil {
		t.Error(err)
	}

	// Check that the address was written into the db
	if _, err := w.GetECAddress(e.PubString()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

func TestGenerateFCTAddress(t *testing.T) {
	// create a new database
	w, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// Generate a new fct address
	f, err := w.GenerateFCTAddress()
	if err != nil {
		t.Error(err)
	}

	// Check that the address was written into the db
	if _, err := w.GetFCTAddress(f.String()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

func TestGetAllAddresses(t *testing.T) {
	e1Sec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	f1Sec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	e2Sec := "Es4NQHwo8F4Z4oMnVwndtjV1rzZN3t5pP5u5jtdgiR1RA6FH4Tmc"
	f2Sec := "Fs3GFV6GNV6ar4b8eGcQWpGFbFtkNWKfEPdbywmha8ez5p7XMJyk"
	correctLen := 2

	// create a new database
	w, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// write a new ec address to the db
	e1, err := factom.GetECAddress(e1Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertECAddress(e1); err != nil {
		t.Error(err)
	}
	e2, err := factom.GetECAddress(e2Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertECAddress(e2); err != nil {
		t.Error(err)
	}

	// write a new fct address to the db
	f1, err := factom.GetFactoidAddress(f1Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertFCTAddress(f1); err != nil {
		t.Error(err)
	}
	f2, err := factom.GetFactoidAddress(f2Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.InsertFCTAddress(f2); err != nil {
		t.Error(err)
	}

	// get all addresses out of db
	fs, es, err := w.GetAllAddresses()
	if err != nil {
		t.Error(err)
	} else if fs == nil {
		t.Errorf("No Factoid address was retrived")
	} else if es == nil {
		t.Errorf("No EC address was retrived")
	}
	// check that all the addresses are there
	if len(fs) != correctLen {
		t.Errorf("Wrong number of factoid addesses were retrived: %v", fs)
	}
	if len(es) != correctLen {
		t.Errorf("Wrong number of ec addesses were retrived: %v", es)
	}

	// print the addresses
	for _, f := range fs {
		t.Logf("%s %s", f, f.SecString())
	}
	for _, e := range es {
		t.Logf("%s %s", e, e.SecString())
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}
