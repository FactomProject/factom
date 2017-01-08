// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"fmt"

	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factoid/state/stateinit"
	"github.com/FactomProject/factoid/wallet"
	"github.com/FactomProject/factom"
)

// This file is a dirty hack to to get the keys out of a version 1 wallet.

// ImportV1Wallet takes a version 1 wallet bolt.db file and imports all of its
// addresses into a factom wallet.
func ImportV1Wallet(v1path, v2path string) (*Wallet, error) {
	w, err := NewOrOpenBoltDBWallet(v2path)
	if err != nil {
		return nil, err
	}

	fstate := stateinit.NewFactoidState(v1path)

	_, values := fstate.GetDB().GetKeysValues([]byte(factoid.W_NAME))
	for _, v := range values {
		we, ok := v.(wallet.IWalletEntry)
		if !ok {
			w.Close()
			return nil, fmt.Errorf("Cannot retrieve addresses from version 1 database")
		}

		switch we.GetType() {
		case "fct":
			f, err := factom.MakeFactoidAddress(we.GetPrivKey(0)[:32])
			if err != nil {
				w.Close()
				return nil, err
			}
			if err := w.InsertFCTAddress(f); err != nil {
				w.Close()
				return nil, err
			}
		case "ec":
			e, err := factom.MakeECAddress(we.GetPrivKey(0)[:32])
			if err != nil {
				w.Close()
				return nil, err
			}
			if err := w.InsertECAddress(e); err != nil {
				w.Close()
				return nil, err
			}
		default:
			return nil, fmt.Errorf("version 1 database returned unknown address type %s %#v", we.GetType(), we)
		}
	}
	return w, err
}
