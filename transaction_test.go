// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	. "github.com/FactomProject/factom"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	if err := NewTransaction("b"); err != nil {
		t.Error(err)
	}
	if txs, err := ListTransactions(); err != nil {
		t.Error(err)
	} else {
		for _, v := range txs {
			t.Log(v)
		}
	}
}
