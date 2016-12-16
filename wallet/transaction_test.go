// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	"testing"

	"github.com/FactomProject/factom"
	. "github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factomd/common/primitives"
)

func TestNewTransaction(t *testing.T) {
	// create a new database
	w1, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// create a new transaction
	if err := w1.NewTransaction("test_tx-01"); err != nil {
		t.Error(err)
	}
	if err := w1.NewTransaction("test_tx-02"); err != nil {
		t.Error(err)
	}

	if len(w1.GetTransactions()) != 2 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}

	if err := w1.DeleteTransaction("test_tx-02"); err != nil {
		t.Error(err)
	}

	if len(w1.GetTransactions()) != 1 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
}

func TestAddInput(t *testing.T) {
	zSec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"

	// create a new database
	w1, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// write a new fct address to the db
	f, err := factom.GetFactoidAddress(zSec)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := w1.InsertFCTAddress(f); err != nil {
		t.Error(err)
	}
	// Get the address back out of the db
	adr, err := w1.GetFCTAddress(f.String())
	if err != nil {
		t.Error(err)
	}

	// create a new transaction
	if err := w1.NewTransaction("tx-01"); err != nil {
		t.Error(err)
	}

	if err := w1.AddInput("tx-01", adr.String(), 5); err != nil {
		t.Error(err)
	}
	if len(w1.GetTransactions()) != 1 {
		t.Errorf("wrong number of transactions %v", w1.GetTransactions())
	}
	t.Log(w1.GetTransactions())

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
}

func TestComposeTrasnaction(t *testing.T) {
	f1Sec := "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"
	//	f1Sec := "Fs1KWJrpLdfucvmYwN2nWrwepLn8ercpMbzXshd1g8zyhKXLVLWj"
	f2Sec := "Fs3GFV6GNV6ar4b8eGcQWpGFbFtkNWKfEPdbywmha8ez5p7XMJyk"
	e1Sec := "Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"

	// create a new database
	w1, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	// write a new fct address to the db
	if in, err := factom.GetFactoidAddress(f1Sec); err != nil {
		t.Error(err)
	} else {
		if err := w1.InsertFCTAddress(in); err != nil {
			t.Error(err)
		}
	}

	// Get the address back out of the db
	f1 := factom.NewFactoidAddress()
	if out, err := factom.GetFactoidAddress(f1Sec); err != nil {
		t.Error(err)
	} else {
		if f, err := w1.GetFCTAddress(out.String()); err != nil {
			t.Error(err)
		} else {
			f1 = f
		}
	}

	// setup a factoid address for receving
	f2, err := factom.GetFactoidAddress(f2Sec)
	if err != nil {
		t.Error(err)
	}

	// setup an ec address for receving
	e1, err := factom.GetECAddress(e1Sec)
	if err != nil {
		t.FailNow()
		t.Error(err)
	}

	// create a new transaction
	if err := w1.NewTransaction("tx-01"); err != nil {
		t.Error(err)
	}
	if err := w1.AddInput("tx-01", f1.String(), 5e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddOutput("tx-01", f2.String(), 3e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddECOutput("tx-01", e1.PubString(), 2e8); err != nil {
		t.Error(err)
	}
	if err := w1.AddFee("tx-01", f1.String(), 10000); err != nil {
		t.Error(err)
	}
	if err := w1.SignTransaction("tx-01", true); err != nil {
		t.Error(err)
	}

	if j, err := w1.ComposeTransaction("tx-01"); err != nil {
		t.Error(err)
	} else {
		t.Log(factom.EncodeJSONString(j))
	}

	// close and remove the testing db
	if err := w1.Close(); err != nil {
		t.Error(err)
	}
}

func TestImportComposedTransaction(t *testing.T) {
	transSig := "020158fe8efb78010100afd89f60646f3e8750c550e4582eca5047546ffef89c13a175985e320232" +
		"bacac81cc428afd7c20001ed0da7057f80dfeb596e6c72c4550c7c7694661dfee2c4a6ba1b903b6ec3e201718b" +
		"5edd2914acc2e4677f336c1a32736e5e9bde13663e6413894f57ec272e28015183427204adbc50623d09ea5a76" +
		"947b4c742e5b56d1483483d5d4336ac12872891be0afae50ce916639dd6db6200e3816d8bd73025b79b7af4de11fcd2105"

	transNoSig := "020158feae8aff010100afd89f60646f3e8750c550e4582eca5047546ffef89c13a175985e320232ba" +
		"cac81cc428afd7c20001ed0da7057f80dfeb596e6c72c4550c7c7694661dfee2c4a6ba1b903b6ec3e201718b5e" +
		"dd2914acc2e4677f336c1a32736e5e9bde13663e6413894f57ec272e2800000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	_ = transSig
	w1, err := NewMapDBWallet()
	if err != nil {
		t.Error(err)
	}

	err = w1.ImportComposedTransaction("txImported", transNoSig)
	if err != nil {
		t.Error(err)
	}

	trans := w1.GetTransactions()["txImported"]
	if trans == nil {
		t.Error("Transaction not found")
	}

	ins := trans.GetInputs()
	if len(ins) != 1 {
		t.Error("Transaction only has 1 input")
	}

	inAddr := primitives.ConvertFctAddressToUserStr(ins[0].GetAddress())
	if inAddr != "FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q" {
		t.Error("Input does not match address")
	}

	outs := trans.GetOutputs()
	if len(outs) != 1 {
		t.Error("Transaction only has 1 input")
	}

	outAddr := primitives.ConvertFctAddressToUserStr(outs[0].GetAddress())
	if outAddr != "FA1yvkgzxMigVDU1WRnpHXhAE3e5zQpE6KKyf5EF76Y34TSg6m8X" {
		t.Error("Ouput does not match address")
	}

	sum, err := trans.TotalOutputs()
	if err != nil {
		t.Error(err)
	}

	if sum != 1e8 {
		t.Error("Output amount is incorrect")
	}

	err = trans.ValidateSignatures()
	if err == nil {
		t.Error("Should be an error")
	}

	// With sig
	err = w1.ImportComposedTransaction("txImportedSig", transSig)
	if err != nil {
		t.Error(err)
	}

	transStructSig := w1.GetTransactions()["txImportedSig"]
	if transStructSig == nil {
		t.Error("Transaction not found")
	}

	err = transStructSig.ValidateSignatures()
	if err != nil {
		t.Error(err)
	}
}
