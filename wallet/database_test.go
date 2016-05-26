// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"os"
	"testing"
)

func TestNewWalletDB(t *testing.T) {
	dbpath := os.TempDir() + "/ldb1"
	
	// create a new database
	w1, err := NewWalletDB(dbpath)
	if err != nil {
		t.Error(err)
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