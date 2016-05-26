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

func TestNewWalletDB(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w1, err := NewWalletDB(dbpath)
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
	if _, err := NewWalletDB(dbpath); err == nil {
		t.Errorf("NewWalletDB did not report error on existing path: %s", dbpath)
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
	w, err := NewWalletDB(dbpath)
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
	w, err := NewWalletDB(dbpath)
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
	w, err := NewWalletDB(dbpath)
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
	w, err := NewWalletDB(dbpath)
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
