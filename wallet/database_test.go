// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"os"
	"testing"

	"github.com/FactomProject/factom"
)

func TestNewWallet(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// check that the seed got written
	w1.lock.RLock()
	seed, err := w1.ldb.Get(seedDBKey, nil)
	if err != nil {
		t.Error(err)
	}
	w1.lock.RUnlock()
	if len(seed) != 64 {
		t.Errorf("stored db seed is the wrong length: %x", seed)
	}
	if bytes.Equal(seed, make([]byte, 64)) {
		t.Errorf("stored db seed is blank")
	}
	
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
	
	// try and create a new database where one already exists 
	if _, err := NewWallet(dbpath); err == nil {
		t.Errorf("NewWallet did not report error on existing path: %s", dbpath)
	}
	
	// remove the testing db
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestOpenWallet(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w1, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	w1.Close()
	
	// make sure we can open the db
	w2, err := OpenWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// check that the seed is there
	w1.lock.RLock()
	seed, err := w1.ldb.Get(seedDBKey, nil)
	if err != nil {
		t.Error(err)
	}
	w1.lock.RUnlock()
	if len(seed) != 64 {
		t.Errorf("stored db seed is the wrong length: %x", seed)
	}
	if bytes.Equal(seed, make([]byte, 64)) {
		t.Errorf("stored db seed is blank")
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
	
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// write a new ec address to the db
	e, err := factom.GetECAddress(zSec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutECAddress(e); err != nil {
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
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestPutFCTAddress(t *testing.T) {
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// write a new fct address to the db
	f, err := factom.GetFactoidAddress(zSec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutFCTAddress(f); err != nil {
		t.Error(err)
	}
	
	// Check that the address was written into the db
	if _, err := w.GetFCTAddress(f.PubString()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestGenerateECAddress(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w, err := NewWallet(dbpath)
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
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestGenerateFCTAddress(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// Generate a new fct address
	f, err := w.GenerateFCTAddress()
	if err != nil {
		t.Error(err)
	}
	
	// Check that the address was written into the db
	if _, err := w.GetFCTAddress(f.PubString()); err != nil {
		t.Error(err)
	}

	// close and remove the testing db
	if err := w.Close(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}

func TestGetAllAddresses(t *testing.T) {
	e1Sec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
	f1Sec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	e2Sec := "Es4NQHwo8F4Z4oMnVwndtjV1rzZN3t5pP5u5jtdgiR1RA6FH4Tmc"
	f2Sec := "Fs3GFV6GNV6ar4b8eGcQWpGFbFtkNWKfEPdbywmha8ez5p7XMJyk"
	correctLen := 2
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w, err := NewWallet(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	// write a new ec address to the db
	e1, err := factom.GetECAddress(e1Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutECAddress(e1); err != nil {
		t.Error(err)
	}
	e2, err := factom.GetECAddress(e2Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutECAddress(e2); err != nil {
		t.Error(err)
	}
	
	// write a new fct address to the db
	f1, err := factom.GetFactoidAddress(f1Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutFCTAddress(f1); err != nil {
		t.Error(err)
	}
	f2, err := factom.GetFactoidAddress(f2Sec)
	if err != nil {
		t.Error(err)
	}
	if err := w.PutFCTAddress(f2); err != nil {
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
	if err := os.RemoveAll(dbpath); err != nil {
		t.Error(err)
	}
}
