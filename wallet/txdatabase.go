// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"fmt"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/hybridDB"
)

// Database keys and key prefixes
var (
	fblockDBPrefix = []byte("FBlock")
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

func (db *TXDatabaseOverlay) Close() error {
	return db.dbo.Close()
}

// GetAllTXs returns a list of all transactions in the history of Factom. A
// local database is used to cache the factoid blocks.
func (db *TXDatabaseOverlay) GetAllTXs() ([]interfaces.ITransaction, error) {
	// update the database and get the newest fblock
	newest, err := db.update()
	if err != nil {
		return nil, err
	}
	fblock, err := db.GetFBlock(newest)
	if err != nil {
		return nil, err
	}
	txs := make([]interfaces.ITransaction, 0)

	for {
		// get all of the txs from the block
		height := fblock.GetDatabaseHeight()
		for _, tx := range fblock.GetTransactions() {
			ins, err := tx.TotalInputs()
			if err != nil {
				return nil, err
			}
			outs, err := tx.TotalOutputs()
			if err != nil {
				return nil, err
			}

			if ins != 0 || outs != 0 {
				tx.SetBlockHeight(height)
				txs = append(txs, tx)
			}
		}

		if pre := fblock.GetPrevKeyMR().String(); pre != factom.ZeroHash {
			// get the previous block
			fblock, err = db.GetFBlock(pre)
			if err != nil {
				return nil, err
			} else if fblock == nil {
				return nil, fmt.Errorf("Missing fblock in database: %s", pre)
			}
		} else {
			break
		}
	}

	return txs, nil
}

// GetTX gets a transaction by the transaction id
func (db *TXDatabaseOverlay) GetTX(txid string) (
	interfaces.ITransaction, error) {
	txs, err := db.GetAllTXs()
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		if tx.GetSigHash().String() == txid {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("Transaction not found")
}

// GetTXAddress returns a list of all transactions in the history of Factom that
// include a specific address.
func (db *TXDatabaseOverlay) GetTXAddress(adr string) (
	[]interfaces.ITransaction, error) {
	filtered := make([]interfaces.ITransaction, 0)

	txs, err := db.GetAllTXs()
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		for _, in := range tx.GetInputs() {
			if primitives.ConvertFctAddressToUserStr(
				in.GetAddress()) == adr {
				txs = append(filtered, tx)
			}
		}
		for _, out := range tx.GetOutputs() {
			if primitives.ConvertFctAddressToUserStr(
				out.GetAddress()) == adr {
				txs = append(filtered, tx)
			}
		}
	}

	return filtered, nil
}

func (db *TXDatabaseOverlay) GetTXRange(start, end int) (
	[]interfaces.ITransaction, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("Range cannot have negative numbers")
	}
	s, e := uint32(start), uint32(end)

	filtered := make([]interfaces.ITransaction, 0)

	txs, err := db.GetAllTXs()
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		if s <= tx.GetBlockHeight() && tx.GetBlockHeight() <= e {
			filtered = append(filtered, tx)
		}
	}

	return filtered, nil
}

// GetFBlock retrives a Factoid Block from Factom
func (db *TXDatabaseOverlay) GetFBlock(keymr string) (interfaces.IFBlock, error) {
	fblock := new(factoid.FBlock)
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

// update gets all fblocks written since the database was last updated, and
// returns the most recent fblock keymr.
func (db *TXDatabaseOverlay) update() (string, error) {
	fblock, err := fblockHead()
	if err != nil {
		return "", err
	}
	newest := fblock.GetKeyMR().String()

	prevmr := fblock.GetPrevKeyMR().String()
	for prevmr != factom.ZeroHash {
		//		// stop when we reach an fblock that is already in the db
		//		if f, err := db.GetFBlock(prevmr); err != nil {
		//			return "", err
		//		} else if f != nil {
		//			return newest, nil
		//		}

		// add the fblock to the db
		if err := db.InsertFBlock(fblock); err != nil {
			return "", err
		}

		fblock, err = getfblock(prevmr)
		if err != nil {
			return "", err
		}
		prevmr = fblock.GetPrevKeyMR().String()
	}

	// write the last fblock into the db
	if err := db.InsertFBlock(fblock); err != nil {
		return "", err
	}

	return newest, nil
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

	return getfblock(fblockmr)
}

func getfblock(keymr string) (interfaces.IFBlock, error) {
	p, err := factom.GetRaw(keymr)
	if err != nil {
		return nil, err
	}
	return factoid.UnmarshalFBlock(p)
}
