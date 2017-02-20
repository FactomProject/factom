// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/FactomProject/factom"
)

func TestUnmarshalJSON(t *testing.T) {
	jsonentry1 := []byte(`
	{
		"ChainID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"ExtIDs":[
			"bbbb",
			"cccc"
		],
		"Content":"111111111111111111"
	}`)

	jsonentry2 := []byte(`
	{
		"ChainName":["aaaa", "bbbb"],
		"ExtIDs":[
			"cccc",
			"dddd"
		],
		"Content":"111111111111111111"
	}`)

	e1 := new(Entry)
	if err := e1.UnmarshalJSON(jsonentry1); err != nil {
		t.Error(err)
	}

	e2 := new(Entry)
	if err := e2.UnmarshalJSON(jsonentry2); err != nil {
		t.Error(err)
	}
}

func TestEntryPrinting(t *testing.T) {
	ent := new(Entry)
	ent.ChainID = "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	ent.Content = []byte("This is a test Entry.")
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the first extid."))
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the second extid."))

	//fmt.Println(ent.String())
	expectedReturn := `EntryHash: 52385948ea3ab6fd67b07664ac6a30ae5f6afa94427a547c142517beaa9054d0
ChainID: 5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069
ExtID: This is the first extid.
ExtID: This is the second extid.
Content:
This is a test Entry.
`

	if ent.String() != expectedReturn {
		fmt.Println(ent.String())
		fmt.Println(expectedReturn)
		t.Fail()
	}

	expectedReturn = `{"chainid":"5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069","extids":["54686973206973207468652066697273742065787469642e","5468697320697320746865207365636f6e642065787469642e"],"content":"546869732069732061207465737420456e7472792e"}`
	jsonReturn, _ := ent.MarshalJSON()
	if string(jsonReturn) != expectedReturn {
		fmt.Println(string(jsonReturn))
		fmt.Println(expectedReturn)
		t.Fail()
	}
}

func TestMarshalBinary(t *testing.T) {
	ent := new(Entry)
	ent.ChainID = "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	ent.Content = []byte("This is a test Entry.")
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the first extid."))
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the second extid."))

	expected, _ := hex.DecodeString("005a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be7090690035001854686973206973207468652066697273742065787469642e00195468697320697320746865207365636f6e642065787469642e546869732069732061207465737420456e7472792e")

	result, _ := ent.MarshalBinary()
	//fmt.Printf("%x\n",result)
	if !bytes.Equal(result, expected) {
		fmt.Printf("found %x expected %x\n", result, expected)
		t.Fail()
	}
}
