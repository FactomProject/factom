// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	. "github.com/FactomProject/factom"

	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factom/wallet/wsapi"
)

func TestImportAddress(t *testing.T) {
	var (
		fs1 = "Fs2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeXA"
		fa1 = "FA3T1gTkuKGG2MWpAkskSoTnfjxZDKVaAYwziNTC1pAYH5B9A1rh"
		es1 = "Es4KmwK65t9HCsibYzVDFrijvkgTFZKdEaEAgfMtYTPSVtM3NDSx"
		ec1 = "EC2CyGKaNddLFxrjkFgiaRZnk77b8iQia3Zj6h5fxFReAcDwCo3i"

		bads = []string{
			"Fs2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeX",  //short
			"Fs2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeeA", //check
			"", //empty
			"Fc2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeXA", //prefix
		}
	)

	// start the test wallet
	done, err := StartTestWallet()
	if err != nil {
		t.Error(err)
	}
	defer func() { done <- 1 }()

	// import the good addresses
	if _, _, err := ImportAddresses(fs1, es1); err != nil {
		t.Error(err)
	}

	if f, err := FetchFactoidAddress(fa1); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("Wallet returned nil factoid address")
	} else if f.SecString() != fs1 {
		t.Error("Wallet returned incorrect address", fs1, f.SecString())
	}

	if e, err := FetchECAddress(ec1); err != nil {
		t.Error(err)
	} else if e == nil {
		t.Error("Wallet returned nil ec address")
	} else if e.SecString() != es1 {
		t.Error("Wallet returned incorrect address", es1, e.SecString())
	}

	// try to import the bad addresses
	for _, bad := range bads {
		if _, _, err := ImportAddresses(bad); err == nil {
			t.Error("Bad address was imported without error", bad)
		}
	}
}

func TestHandleWalletBalances(t *testing.T) {
	// start the test wallet
	done, err := StartTestWallet()
	if err != nil {
		t.Error(err)
	}
	defer func() { done <- 1 }()

	// Testing when all accounts dont have balances #2
	noBalFCT := "Fs1itDLe8GoFCLsdbqb2rs6U67wQX4TikkTJV69BxGuG1tDvs41q"
	noBalEC := "Es3W3R2u85aN2MNr2EoyMazAy7yGTZZg8eDaW7vfjorkrWAANv6t"

	addr2 := []string{noBalFCT, noBalEC}
	testingVar2, _ := helper(t, addr2)
	if testingVar2.Result.FactoidAccountBalances.Ack != 0 && testingVar2.Result.FactoidAccountBalances.Saved != 0 && testingVar2.Result.EntryCreditAccountBalances.Ack != 0 && testingVar2.Result.EntryCreditAccountBalances.Saved != 0 {
		t.Error("balances are not what they should be")
	}
	fmt.Println("Passed balance of 0 #2")

	// Testing when all accounts have balances #3
	hasBalFCT := "Fs1vEcszU16mC72CBMAfAnxVvKQKTtrTqiCfdGF8hycMn1j1DBKy"
	hasBalEC := "Es2nSXmiaUuk9AxX2X43Ws4XjXPCxehTyHZAEn5NJH9ei1gLW1FR"

	addr3 := []string{hasBalFCT, hasBalEC}
	testingVar3, _ := helper(t, addr3)
	if testingVar3.Result.EntryCreditAccountBalances.Ack != 40 && testingVar3.Result.EntryCreditAccountBalances.Saved != 40 && testingVar3.Result.FactoidAccountBalances.Ack != 0 && testingVar3.Result.FactoidAccountBalances.Saved != 0 {
		t.Error("balances are not what they should be")
	}
	fmt.Println("Passed when some have values #3")
}

type walletcall struct {
	Jsonrpc string `json:"jsonrps"`
	Id      int    `json:"id"`
	Result  struct {
		FactoidAccountBalances struct {
			Ack   int64 `json:"ack"`
			Saved int64 `json:"saved"`
		} `json:"fctaccountbalances"`
		EntryCreditAccountBalances struct {
			Ack   int64 `json:"ack"`
			Saved int64 `json:"saved"`
		} `json:"ecaccountbalances"`
	} `json:"result"`
}

func helper(t *testing.T, addr []string) (*walletcall, string) {
	for _, k := range addr {
		if _, _, err := ImportAddresses(k); err != nil {
			return nil, "failed"
		}
	}

	url := "http://localhost:8089/v2"
	jsonStrEC := []byte(`{"jsonrpc": "2.0", "id": 0, "method": "wallet-balances"}`)
	reqEC, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStrEC))
	reqEC.Header.Set("content-type", "text/plain;")

	clientEC := &http.Client{}
	callRespEC, err := clientEC.Do(reqEC)
	if err != nil {
		t.Error(err)
	}

	defer callRespEC.Body.Close()
	bodyEC, _ := ioutil.ReadAll(callRespEC.Body)
	fmt.Println("BODY: ", string(bodyEC))

	respEC := new(walletcall)
	errEC := json.Unmarshal([]byte(bodyEC), &respEC)
	if errEC != nil {
		t.Error(errEC)
	}
	return respEC, ""
}

func TestImportKoinify(t *testing.T) {
	var (
		good_mnemonic = "yellow yellow yellow yellow yellow yellow yellow" +
			" yellow yellow yellow yellow yellow" // good
		koinifyexpect = "FA3cih2o2tjEUsnnFR4jX1tQXPpSXFwsp3rhVp6odL5PNCHWvZV1"

		bad_mnemonic = []string{
			"", // bad empty
			"yellow yellow yellow yellow yellow yellow yellow yellow yellow" +
				" yellow yellow", // bad short
			"yellow yellow yellow yellow yellow yellow yellow yellow yellow" +
				" yellow yellow asdfasdf", // bad word
		}
	)

	// start the test wallet
	done, err := StartTestWallet()
	if err != nil {
		t.Error(err)
	}
	defer func() { done <- 1 }()

	// check the import for koinify names
	fa, err := ImportKoinify(good_mnemonic)
	if err != nil {
		t.Error(err)
	}
	if fa.String() != koinifyexpect {
		t.Error("Incorrect address from Koinify mnemonic", fa, koinifyexpect)
	}

	for _, m := range bad_mnemonic {
		if _, err := ImportKoinify(m); err == nil {
			t.Error("No error for bad address:", m)
		}
	}
}

// helper functions for testing

func populateTestWallet() error {
	//FA3T1gTkuKGG2MWpAkskSoTnfjxZDKVaAYwziNTC1pAYH5B9A1rh
	//Fs2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeXA
	//
	//FA3oaS2D2GkrZJuWuiDohnLruxV3AWbrM3PmG3HSSE7DHzPWio36
	//Fs1os7xg2mN9fTuJmaYZLk6EXz51x2wmmHr2365UAuPMJW3aNr25
	//
	//EC2CyGKaNddLFxrjkFgiaRZnk77b8iQia3Zj6h5fxFReAcDwCo3i
	//Es4KmwK65t9HCsibYzVDFrijvkgTFZKdEaEAgfMtYTPSVtM3NDSx
	//
	//EC2R4bPDj9WQ8eWA4X3K8NYfTkBh4HFvCopLBq48FyrNXNumSK6w
	//Es355qB6tWo1ZZRTK8cXpHjxGECXaPGw98AFCRJ6kxZ3J6vp1M2i

	_, _, err := ImportAddresses(
		"Fs2TCa7Mo4XGy9FQSoZS8JPnDfv7SjwUSGqrjMWvc1RJ9sKbJeXA",
		"Fs1os7xg2mN9fTuJmaYZLk6EXz51x2wmmHr2365UAuPMJW3aNr25",
		"Es4KmwK65t9HCsibYzVDFrijvkgTFZKdEaEAgfMtYTPSVtM3NDSx",
		"Es355qB6tWo1ZZRTK8cXpHjxGECXaPGw98AFCRJ6kxZ3J6vp1M2i",
	)
	if err != nil {
		return err
	}

	return nil
}

// StartTestWallet runs a test wallet and serves the wallet api. The caller
// must write an int to the chan when compleate to stop the wallet api and
// remove the test db.
func StartTestWallet() (chan int, error) {
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
