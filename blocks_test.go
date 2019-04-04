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

func TestEblockString(t *testing.T) {
	b := new(EBlock)

	b.Header.BlockSequenceNumber = 50
	b.Header.ChainID = "6909765ff072c322c56a7c4bfa8911ee4fdefacca711d30a9ad2a8672a3cc959"
	b.Header.PrevKeyMR = "5d94cc642a9cccdc61a8926b7ddc1223dfe26a1ffdc71f597b8b22fc73a8a3a0"
	b.Header.Timestamp = 1487615370
	b.Header.DBHeight = 76802

	e := EBEntry{EntryHash: "125c9be87883666f5a0afa22424328a7af8df3aa3dd6984890bb096c1a8a11ae", Timestamp: 1487615370}
	b.EntryList = append(b.EntryList, e)
	//fmt.Println(b)
	expectedEntryString := `BlockSequenceNumber: 50
ChainID: 6909765ff072c322c56a7c4bfa8911ee4fdefacca711d30a9ad2a8672a3cc959
PrevKeyMR: 5d94cc642a9cccdc61a8926b7ddc1223dfe26a1ffdc71f597b8b22fc73a8a3a0
Timestamp: 1487615370
DBHeight: 76802
EBEntry {
	Timestamp 1487615370
	EntryHash 125c9be87883666f5a0afa22424328a7af8df3aa3dd6984890bb096c1a8a11ae
}
`
	if b.String() != expectedEntryString {
		fmt.Println(b.String())
		fmt.Println(expectedEntryString)
		t.Fail()
	}
}
