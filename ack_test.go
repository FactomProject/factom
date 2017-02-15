// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"
	"fmt"
	
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


