// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"fmt"
	"testing"

	. "github.com/FactomProject/factom"
)

var ()

type doubleString struct {
	ChainID string `json:"chainid"`
	KeyMR   string `json:"keymr"`
}

func TestDblockString(t *testing.T) {
	d := new(DBlock)

	d.Header.PrevBlockKeyMR = "fc01feda95f64a697431de6283593012a299cee7e834061a62b4addf0756dc2d"
	d.Header.Timestamp = 1487615370
	d.Header.SequenceNumber = 76802

	e := doubleString{ChainID: "6909765ff072c322c56a7c4bfa8911ee4fdefacca711d30a9ad2a8672a3cc959", KeyMR: "3241ef7e4122a5f8c9df4536370e0a8919d5d20593fee0d49459e67586e56742"}
	d.EntryBlockList = append(d.EntryBlockList, e)
	//fmt.Println(d)
	expectedEntryString := `PrevBlockKeyMR: fc01feda95f64a697431de6283593012a299cee7e834061a62b4addf0756dc2d
Timestamp: 1487615370
SequenceNumber: 76802
EntryBlock {
	ChainID 6909765ff072c322c56a7c4bfa8911ee4fdefacca711d30a9ad2a8672a3cc959
	KeyMR 3241ef7e4122a5f8c9df4536370e0a8919d5d20593fee0d49459e67586e56742
}
`
	if d.String() != expectedEntryString {
		fmt.Println(d.String())
		fmt.Println(expectedEntryString)
		t.Fail()
	}
}

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
