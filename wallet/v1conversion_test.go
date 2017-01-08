// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	"os"
	"testing"

	. "github.com/FactomProject/factom/wallet"
)

func TestImportV1Wallet(t *testing.T) {
	v1path := os.TempDir() + "/factoid_wallet_bolt.db"
	v2path := os.TempDir() + "/test_wallet-01"

	w, err := ImportV1Wallet(v1path, v2path)
	if err != nil {
		t.Error(err)
	}

	fs, es, err := w.GetAllAddresses()
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
	if err := os.RemoveAll(v2path); err != nil {
		t.Error(err)
	}
}
