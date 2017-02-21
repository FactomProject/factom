// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestComposeEntryCommit(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}
	ecAddr, _ := GetECAddress("Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG")
	ent := new(Entry)
	ent.ChainID = "954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f4"
	ent.Content = []byte("test!")
	ent.ExtIDs = append(ent.ExtIDs, []byte("test"))

	eCommit, _ := ComposeEntryCommit(ent, ecAddr)
	r := new(response)
	json.Unmarshal(eCommit.Params, r)
	binCommit, _ := hex.DecodeString(r.Message)

	//fmt.Printf("%x\n",binCommit)
	//the commit has a timestamp which is updated new for each time it is called.  This means it is different after each call.
	//we will check the non-changing parts

	if len(binCommit) != 136 {
		fmt.Println("expected commit to be 136 bytes long, instead got", len(binCommit))
		t.Fail()
	}
	result := binCommit[0:1]
	expected := []byte{0x00}
	if !bytes.Equal(result, expected) {
		fmt.Printf("found %x expected %x\n", result, expected)
		t.Fail()
	}
	//skip the 6 bytes of the timestamp
	result = binCommit[7:72]
	expected, _ = hex.DecodeString("285ED45081D5B8819A678D13C7C2D04F704B34C74E8AAECD9BD34609BEE04720013B6A27BCCEB6A42D62A3A8D02A6F0D73653215771DE243A63AC048A18B59DA29")

	if !bytes.Equal(result, expected) {
		fmt.Printf("found %x expected %x\n", result, expected)
		t.Fail()
	}
}

func TestComposeEntryReveal(t *testing.T) {

	ent := new(Entry)
	ent.ChainID = "954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f4"
	ent.Content = []byte("test!")
	ent.ExtIDs = append(ent.ExtIDs, []byte("test"))

	eReveal, _ := ComposeEntryReveal(ent)

	expectedResponse := `{"entry":"00954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f400060004746573747465737421"}`
	if expectedResponse != string(eReveal.Params) {
		fmt.Println(eReveal.Params)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestCommitEntry(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "message": "Entry Commit Success",
    "txid": "bf12150038699f678ac2314e9fa2d4786dc8984d9b8c67dab8cd7c2f2e83372c"
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	ent := new(Entry)
	ent.ChainID = "954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f4"
	ent.Content = []byte("test!")
	ent.ExtIDs = append(ent.ExtIDs, []byte("test"))
	ecAddr, _ := GetECAddress("Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG")

	response, _ := CommitEntry(ent, ecAddr)

	//fmt.Println(response)
	expectedResponse := "bf12150038699f678ac2314e9fa2d4786dc8984d9b8c67dab8cd7c2f2e83372c"

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestReveaEntry(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "message": "Entry Reveal Success",
    "entryhash": "f5c956749fc3eba4acc60fd485fb100e601070a44fcce54ff358d60669854734"
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	ent := new(Entry)
	ent.ChainID = "954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f4"
	ent.Content = []byte("test!")
	ent.ExtIDs = append(ent.ExtIDs, []byte("test"))

	response, _ := RevealEntry(ent)

	//fmt.Println(response)
	expectedResponse := "f5c956749fc3eba4acc60fd485fb100e601070a44fcce54ff358d60669854734"

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}
