// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	. "github.com/FactomProject/factom/wallet"
	
	"os"
	"testing"
	
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
)

func TestTXDatabaseOverlay(t *testing.T) {
	dbpath := os.TempDir() + "/test_txdb-01"
	db1, err := NewTXLevelDB(dbpath)
	if err != nil {
		t.Error(err)
	}
	
	fblock, err := fblockHead()
	if err != nil {
		t.Error(err)
	}	
	if err := db1.InsertFBlock(fblock); err != nil {
		t.Error(err)
	}	
	if f, err := db1.GetFBlock(fblock.GetKeyMR().String()); err != nil {
		t.Error(err)
	}
}

// fblockHead gets the most recent fblock.
func fblockHead() (interfaces.IFBlock, error) {
		fblockID := "000000000000000000000000000000000000000000000000000000000000000f"

		dbhead, err := factom.GetDBlockHead()
		if err != nil {
			return nil, err
		}
		dblock, err := factom.GetDBlock(dbhead)
		if err != nil {
			return nil, err
		}
		
		var fblockmr string
		for _, eblock := range dblock.EntryBlockList {
			if eblock.ChainID == fblockID {
				fblockmr = eblock.KeyMR
			}
		}
		if fblockmr == "" {
			return nil, err
		}
		
		// get the most recent block
		p, err := factom.GetRaw(fblockmr)
		if err != nil {
			return nil, err
		}
		return factoid.UnmarshalFBlock(p)
}
