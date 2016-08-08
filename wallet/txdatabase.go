// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"fmt"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/hybridDB"
)

// Database keys and key prefixes
var (
	fblockDBPrefix    = []byte("FBlock")
)

type TXDatabaseOverlay struct {
	dbo databaseOverlay.Overlay
}

func NewTXOverlay(db interfaces.IDatabase) *TXDatabaseOverlay {
	answer := new(TXDatabaseOverlay)
	answer.dbo.DB = db
	return answer
}

func NewTXLevelDB(ldbpath string) (*TXDatabaseOverlay, error) {
	db, err := hybridDB.NewLevelMapHybridDB(ldbpath, false)
	if err != nil {
		fmt.Printf("err opening transaction db: %v\n", err)
	}

	if db == nil {
		fmt.Println("Creating new transaction db ...")
		db, err = hybridDB.NewLevelMapHybridDB(ldbpath, true)

		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Transaction database started from: " + ldbpath)
	return NewTXOverlay(db), nil
}

func (db *TXDatabaseOverlay) GetFBlock(keymr string) (interfaces.IFBlock, error) {
	var fblock interfaces.IFBlock
	data, err := db.dbo.Get(fblockDBPrefix, []byte(keymr), fblock)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data.(interfaces.IFBlock), nil
}

func (db *TXDatabaseOverlay) InsertFBlock(fblock interfaces.IFBlock) error {
	fblockmr := []byte(fblock.GetKeyMR().String())

	batch := []interfaces.Record{}
	batch = append(batch, interfaces.Record{fblockDBPrefix, fblockmr, fblock})

	return db.dbo.PutInBatch(batch)
}

// update gets all fblocks written since the database was last updated.
func (db *TXDatabaseOverlay) update() error {
	return nil
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

//var getAllTransactions = func() *fctCmd {
//	cmd := new(fctCmd)
//	cmd.helpMsg = "factom-cli get alltxs"
//	cmd.description = "Get the entire history of transactions in Factom"
//	cmd.execFunc = func(args []string) {
//		fblockID := "000000000000000000000000000000000000000000000000000000000000000f"
//
//		dbhead, err := factom.GetDBlockHead()
//		if err != nil {
//			errorln(err)
//			return
//		}
//		dblock, err := factom.GetDBlock(dbhead)
//		if err != nil {
//			errorln(err)
//			return
//		}
//		
//		var fblockmr string
//		for _, eblock := range dblock.EntryBlockList {
//			if eblock.ChainID == fblockID {
//				fblockmr = eblock.KeyMR
//			}
//		}
//		if fblockmr == "" {
//			errorln("no fblock in current dblock")
//			return
//		}
//		
//		// get the most recent block
//		p, err := factom.GetRaw(fblockmr)
//		if err != nil {
//			errorln(err)
//			return
//		}
//		fblock, err := factoid.UnmarshalFBlock(p)
//		if err != nil {
//			errorln(err)
//			return
//		}
//		
//		for fblock.GetPrevKeyMR().String() != factom.ZeroHash {
//			txs := fblock.GetTransactions()
//			for _, tx := range txs {
//				fmt.Println(tx)
//			}
//			p, err := factom.GetRaw(fblock.GetPrevKeyMR().String())
//			if err != nil {
//				errorln(err)
//				return
//			}
//			fblock, err = factoid.UnmarshalFBlock(p)
//			if err != nil {
//				errorln(err)
//				return
//			}
//		}
//		
//		// print the first fblock
//		txs := fblock.GetTransactions()
//		for _, tx := range txs {
//			fmt.Println(tx)
//		}
//	}
//	help.Add("get alltxs", cmd)
//	return cmd
//}()
