// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	"testing"

	//"github.com/FactomProject/factom"
	. "github.com/FactomProject/factom/wallet"
)

func TestImportWithSpaces(t *testing.T) {
	w, err := ImportWalletFromMnemonic("yellow  yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = w.GenerateFCTAddress()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
