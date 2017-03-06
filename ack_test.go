// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/FactomProject/factom"
)

var ()

func TestAckStrings(t *testing.T) {
	status := new(EntryStatus)
	status.CommitTxID = "107c239ee41bb2b0cfa19d8760deb82c942f1bac8ad99516f2f801bf16ae2998"
	//status.EntryHash = "1b363e01af0c0e28f0acbc33bc816ec11f4b28680797e74e341476409dd52295"
	gtd := new(GeneralTransactionData)
	gtd.Status = "TransactionACK"
	gtd.TransactionDateString = "2017-02-15 13:01:41"

	status.CommitData = *gtd

	entryPrintout := status.String()
	//fmt.Println(entryPrintout)

	expectedString := `TxID: 107c239ee41bb2b0cfa19d8760deb82c942f1bac8ad99516f2f801bf16ae2998
Status: TransactionACK
Date: 2017-02-15 13:01:41
`
	if entryPrintout != expectedString {
		fmt.Println(entryPrintout)
		fmt.Println(expectedString)
		t.Fail()
	}

	txstatus := new(FactoidTxStatus)
	txstatus.TxID = "b8b12fba54bd1857b0262bba1b71dbeb4e17404570c2ebe50de0dabf061d575c"
	//gtdfct := new(GeneralTransactionData)
	txstatus.Status = "TransactionACK"
	txstatus.TransactionDateString = "2017-02-15 15:07:27"
	//txstatus.CommitData = *gtdfct
	fctPrintout := txstatus.String()
	//fmt.Println(fctPrintout)

	expectedfctString := `TxID: b8b12fba54bd1857b0262bba1b71dbeb4e17404570c2ebe50de0dabf061d575c
Status: TransactionACK
Date: 2017-02-15 15:07:27
`
	if fctPrintout != expectedfctString {
		fmt.Println(fctPrintout)
		fmt.Println(expectedfctString)
		t.Fail()
	}
}

func TestAckFct(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{  
   "jsonrpc":"2.0",
   "id":0,
   "result":{  
      "txid":"f1d9919829fa71ce18caf1bd8659cce8a06c0026d3f3fffc61054ebb25ebeaa0",
      "transactiondate":1441138021975,
      "transactiondatestring":"2015-09-01 15:07:01",
      "blockdate":1441137600000,
      "blockdatestring":"2015-09-01 15:00:00",
      "status":"DBlockConfirmed"
   }
}`)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	//fmt.Println("exposed URL:",url)
	SetFactomdServer(url)

	//tx := "02015a43cc6d37010100afd7c200031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6fafd6e4200ceb0a10711f9fb61bc983cb4761817e4b3ff6c31ab0d5da6afb03625e368859013b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29dcc6c027a9d321129381d2d8badb3ccd591fd8a515166ca09a8a72cbf3837916c8e4789b0452dffc708ccde097163a86fd0ac23b11416cebb7ccebcdadbba908"
	//txid := "d998c577a9da5dab3d5634753db3e377e392d72d0204d31bd922df483546da4d"
	tx := "dummy1"
	txid := "dummy2"

	txStatus, _ := FactoidACK(txid, tx)
	//fmt.Println(txStatus)

	expectedfctString := `TxID: f1d9919829fa71ce18caf1bd8659cce8a06c0026d3f3fffc61054ebb25ebeaa0
Status: DBlockConfirmed
Date: 2015-09-01 15:07:01
`
	if txStatus.String() != expectedfctString {
		fmt.Println(txStatus.String())
		fmt.Println(expectedfctString)
		t.Fail()
	}
}

func TestAckEntry(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{  
   "jsonrpc":"2.0",
   "id":0,
   "result":{  
      "committxid":"e5b5be39a41df43a3c46beaa238dc5e6f7bb11115a8da1a9b45cd694e257935a",
      "entryhash":"9228b4b080b3cf94cceea866b74c48319f2093f56bd5a63465288e9a71437ee8",
      "commitdata":{  
         "transactiondate":1449547801861,
         "transactiondatestring":"2015-12-07 22:10:01",
         "blockdate":1449547800000,
         "blockdatestring":"2015-12-07 22:10:00",
         "status":"DBlockConfirmed"
      },
      "entrydata":{  
         "blockdate":1449547800000,
         "blockdatestring":"2015-12-07 22:10:00",
         "status":"DBlockConfirmed"
      }
   }
}`)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	//fmt.Println("exposed URL:",url)
	SetFactomdServer(url)

	tx := "dummy1"
	txid := "dummy2"

	entryStatus, _ := EntryACK(txid, tx)
	//fmt.Println(entryStatus)

	expectedEntryString := `TxID: e5b5be39a41df43a3c46beaa238dc5e6f7bb11115a8da1a9b45cd694e257935a
Status: DBlockConfirmed
Date: 2015-12-07 22:10:01
`
	if entryStatus.String() != expectedEntryString {
		fmt.Println(entryStatus.String())
		fmt.Println(expectedEntryString)
		t.Fail()
	}
}
