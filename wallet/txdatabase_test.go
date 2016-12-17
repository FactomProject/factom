// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	//"fmt" // DEBUG
	"testing"

	. "github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/testHelper"
)

func TestTXDatabaseOverlay(t *testing.T) {
	db1 := NewTXMapDB()

	fblock := fblockHead()
	if err := db1.InsertFBlockHead(fblock); err != nil {
		t.Error(err)
	}
	if f, err := db1.GetFBlock(fblock.GetKeyMR().String()); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Errorf("Fblock not found in db")
	}
}

/*
func TestGetAllTXs(t *testing.T) {
	db1 := NewTXMapDB()

	txs, err := db1.GetAllTXs()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("DEBUG", txs)
	t.Logf("got %d txs", len(txs))
}

func TestGetTXAddress(t *testing.T) {
	db1 := NewTXMapDB()

	adr := "FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"
	txs, err := db1.GetTXAddress(adr)
	if err != nil {
		t.Error(err)
	}
	t.Logf("got %d txs", len(txs))
}
*/
//func TestGetAllTXs(t *testing.T) {
//	dbpath := os.TempDir() + "/test_txdb-01"
//	db1, err := NewTXLevelDB(dbpath)
//	if err != nil {
//		t.Error(err)
//	}
//	defer db1.Close()
//
//	txs := make(chan interfaces.ITransaction, 500)
//	errs := make(chan error)
//	output := make(chan string)
//
//	fmt.Println("DEBUG: running getalltxs")
//	go db1.GetAllTXs(txs, errs)
//
//	go func() {
//		for tx := range txs {
//			output <- fmt.Sprint("Got TX:", tx)
//		}
//		output <- fmt.Sprint("end of txs")
//	}()
//
//	go func() {
//		for err := range errs {
//			output <- fmt.Sprintln("Got error:", err)
//		}
//		output <- fmt.Sprint("end of errs")
//	}()
//
//	for {
//		fmt.Println(<-output)
//	}
//
////	for {
////		select {
////		case tx, ok := <-txs:
////			fmt.Println("Got TX:", tx)
////			if !ok {
////				txs = nil
////			}
////		case err, ok := <-errs:
////			fmt.Println("DEBUG: got error:", err)
////			if !ok {
////				errs = nil
////			}
////		}
////
////		if txs == nil && errs == nil {
////			break
////		}
////	}
//}

// fblockHead gets the most recent fblock.
func fblockHead() interfaces.IFBlock {
	return testHelper.CreateTestFactoidBlock(nil)
}
