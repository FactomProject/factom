// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	. "github.com/FactomProject/factom"
	"testing"

	"encoding/json"
	"os"
	"time"

	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factom/wallet/wsapi"
)

func TestJSONTransactions(t *testing.T) {
	tx1 := mkdummytx()
	t.Log("Transaction:", tx1)
	p, err := json.Marshal(tx1)
	if err != nil {
		t.Error(err)
	}
	t.Log("JSON transaction:", string(p))

	tx2 := new(Transaction)
	if err := json.Unmarshal(p, tx2); err != nil {
		t.Error(err)
	}
	t.Log("Unmarshaled:", tx2)
}

func TestTransactions(t *testing.T) {
	// start the test wallet
	done, err := startTestWallet()
	if err != nil {
		t.Error(err)
	}
	defer func() { done <- 1 }()

	// make sure the wallet is empty
	if txs, err := ListTransactionsTmp(); err != nil {
		t.Error(err)
	} else if len(txs) > 0 {
		t.Error("Unexpected transactions returned from the wallet:", txs)
	}

	// create a new transaction
	tx1, err := NewTransaction("tx1")
	if err != nil {
		t.Error(err)
	}
	if tx1 == nil {
		t.Error("No transaction was returned")
	}

	if tx, err := GetTmpTransaction("tx1"); err != nil {
		t.Error(err)
	} else if tx == nil {
		t.Error("Temporary transaction was not saved in the wallet")
	}

	// delete a transaction
	if err := DeleteTransaction("tx1"); err != nil {
		t.Error(err)
	}

	if txs, err := ListTransactionsTmp(); err != nil {
		t.Error(err)
	} else if len(txs) > 0 {
		t.Error("Unexpected transactions returned from the wallet:", txs)
	}
}

// helper functions for testing

func startTestWallet() (chan int, error) {
	var (
		walletdbfile = os.TempDir() + "/testingwallet.bolt"
		txdbfile     = os.TempDir() + "/testingtxdb.bolt"
	)

	// make a chan to signal when the test is finished with the wallet
	done := make(chan int, 1)

	// setup a testing wallet
	fctWallet, err := wallet.NewOrOpenBoltDBWallet(walletdbfile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(walletdbfile)

	txdb, err := wallet.NewTXBoltDB(txdbfile)
	if err != nil {
		return nil, err
	} else {
		fctWallet.AddTXDB(txdb)
	}
	defer os.Remove(txdbfile)

	RpcConfig = &RPCConfig{
		WalletTLSEnable:   false,
		WalletTLSKeyFile:  "",
		WalletTLSCertFile: "",
		WalletRPCUser:     "",
		WalletRPCPassword: "",
		WalletServer:      "localhost:8089",
	}

	go wsapi.Start(fctWallet, ":8089", *RpcConfig)
	go func() {
		<-done
		wsapi.Stop()
		fctWallet.Close()
		txdb.Close()
	}()

	return done, nil
}

func mkdummytx() *Transaction {
	tx := &Transaction{
		BlockHeight: 42,
		Name:        "dummy",
		Timestamp: func() time.Time {
			t, _ := time.Parse("2006-Jan-02 15:04", "1988-Jan-02 10:00")
			return t
		}(),
		TotalInputs:    13,
		TotalOutputs:   12,
		TotalECOutputs: 1,
	}
	return tx
}
