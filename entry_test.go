// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

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

	expectedReturn := `EntryHash: 52385948ea3ab6fd67b07664ac6a30ae5f6afa94427a547c142517beaa9054d0
ChainID: 5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069
ExtID: This is the first extid.
ExtID: This is the second extid.
Content:
This is a test Entry.
`

	if ent.String() != expectedReturn {
		t.Errorf("expected:%s\nrecieved:%s", expectedReturn, ent.String())
	}
	t.Log(ent.String())

	expectedReturn = `{"chainid":"5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069","extids":["54686973206973207468652066697273742065787469642e","5468697320697320746865207365636f6e642065787469642e"],"content":"546869732069732061207465737420456e7472792e"}`
	jsonReturn, _ := ent.MarshalJSON()
	if string(jsonReturn) != expectedReturn {
		t.Errorf("expected:%s\nrecieved:%s", expectedReturn, string(jsonReturn))
	}
	t.Log(string(jsonReturn))
}

func TestMarshalBinary(t *testing.T) {
	ent := new(Entry)
	ent.ChainID = "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	ent.Content = []byte("This is a test Entry.")
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the first extid."))
	ent.ExtIDs = append(ent.ExtIDs, []byte("This is the second extid."))

	expected, err := hex.DecodeString("005a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be7090690035001854686973206973207468652066697273742065787469642e00195468697320697320746865207365636f6e642065787469642e546869732069732061207465737420456e7472792e")
	if err != nil {
		t.Error(err)
	}

	result, _ := ent.MarshalBinary()
	if !bytes.Equal(result, expected) {
		t.Errorf("expected:%s\nrecieved:%s", expected, result)
	}
	t.Logf("%x", result)
}

func TestNewEntry(t *testing.T) {
	expected, err := hex.DecodeString("005a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be7090690035001854686973206973207468652066697273742065787469642e00195468697320697320746865207365636f6e642065787469642e546869732069732061207465737420456e7472792e")
	if err != nil {
		t.Error(err)
	}

	// Test entry from strings
	efs := NewEntryFromStrings(
		"5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069",
		"This is a test Entry.",
		"This is the first extid.",
		"This is the second extid.",
	)
	efsResult, err := efs.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, efsResult) {
		t.Errorf("expected:%s\nrecieved:%s", expected, efsResult)
	}
	t.Logf("%x", efsResult)

	chainid, err := hex.DecodeString("5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069")
	if err != nil {
		t.Error(err)
	}
	content, err := hex.DecodeString("546869732069732061207465737420456e7472792e")
	if err != nil {
		t.Error(err)
	}
	ext1, err := hex.DecodeString("54686973206973207468652066697273742065787469642e")
	if err != nil {
		t.Error(err)
	}
	ext2, err := hex.DecodeString("5468697320697320746865207365636f6e642065787469642e")
	if err != nil {
		t.Error(err)
	}

	efb := NewEntryFromBytes(
		chainid,
		content,
		ext1,
		ext2,
	)
	efbResult, err := efb.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected, efbResult) {
		t.Errorf("expected:%s\nrecieved:%s", expected, efbResult)
	}
	t.Logf("%x", efbResult)
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

	//the commit has a timestamp which is updated new for each time it is called.  This means it is different after each call.
	//we will check the non-changing parts

	if len(binCommit) != 136 {
		t.Error("expected commit to be 136 bytes long, instead got", len(binCommit))
	}
	result := binCommit[0:1]
	expected := []byte{0x00}
	if !bytes.Equal(result, expected) {
		t.Errorf("expected:%s\nrecieved:%s", expected, result)
	}
	//skip the 6 bytes of the timestamp
	result = binCommit[7:72]
	expected, _ = hex.DecodeString("285ED45081D5B8819A678D13C7C2D04F704B34C74E8AAECD9BD34609BEE04720013B6A27BCCEB6A42D62A3A8D02A6F0D73653215771DE243A63AC048A18B59DA29")

	if !bytes.Equal(result, expected) {
		t.Errorf("expected:%s\nrecieved:%s", expected, result)
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
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, eReveal.Params)
	}
	t.Log(string(eReveal.Params))
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

	expectedResponse := "bf12150038699f678ac2314e9fa2d4786dc8984d9b8c67dab8cd7c2f2e83372c"

	if expectedResponse != response {
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}

func TestReveaEntry(t *testing.T) {
	factomdResponse := `{
	  "jsonrpc": "2.0",
	  "id": 0,
	  "result": {
	    "message": "Entry Reveal Success",
	    "entryhash": "f5c956749fc3eba4acc60fd485fb100e601070a44fcce54ff358d60669854734"
	  }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	ent := new(Entry)
	ent.ChainID = "954d5a49fd70d9b8bcdb35d252267829957f7ef7fa6c74f88419bdc5e82209f4"
	ent.Content = []byte("test!")
	ent.ExtIDs = append(ent.ExtIDs, []byte("test"))

	response, err := RevealEntry(ent)
	if err != nil {
		t.Error(err)
	}

	expectedResponse := "f5c956749fc3eba4acc60fd485fb100e601070a44fcce54ff358d60669854734"

	if expectedResponse != response {
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}

func TestGetEntry(t *testing.T) {
	factomdResponse := `{
		"jsonrpc":"2.0",
		"id":0,
		"result":{
			"chainid":"df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604",
			"content":"68656C6C6F20776F726C64",
			"extids":[
				"466163746f6d416e63686f72436861696e"
			]
		}
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetEntry("be5216cc7a5a3ad44b49245aec298f47cbdfca9862dee13b0093e5880012b771")
	if err != nil {
		t.Error(err)
	}

	expectedResponse := `EntryHash: 1c840bc18be182e89e12f9e63fb8897d13b071b631ced7e656837ccea8fdb3ae
ChainID: df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
ExtID: FactomAnchorChain
Content:
hello world
`

	if expectedResponse != response.String() {
		t.Errorf("expected:%s\nreceived:%s", expectedResponse, response)
	}
}

func TestGetPending(t *testing.T) {
	factomdResponse := `{
		"jsonrpc":"2.0",
		"id":0,
		"result":[
			{"entryhash":"64ebbb592737e8dcad8bce5db209713a8af18ca306ea6542e86bb1d1c88f5529","chainid":"cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d","status":"TransactionACK"},
			{"entryhash":"fcfafbef8471d99cfda1cb8d77a15d71d39d96acaf85da836c1f5075a71bb7c3","chainid":"cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d","status":"TransactionACK"},
			{"entryhash":"5ecfd88edacbaeea8dde7ead5634dbfb000c37e9d6a07ffa14432f3a1d07549d","chainid":"6e4540d08d5ac6a1a394e982fb6a2ab8b516ee751c37420055141b94fe070bfe","status":"TransactionACK"},
			{"entryhash":"d72101ca339aa7206ae6d82594b61051ff901f78276c389ce8b885f34caa9d34","chainid":null,"status":"TransactionACK"},
			{"entryhash":"2eb17336caf62f1eb869a5380298ced3b2653b496bd90e99a0a2da58c8b3101b","chainid":"a642a8674f46696cc47fdb6b65f9c87b2a19c5ea8123b3d2f0c13b6f33a9d5ef","status":"TransactionACK"}
	]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetPendingEntries()
	if err != nil {
		t.Error(err)
	}

	received := fmt.Sprintf("%+v", response)
	expected := "[{ChainID:cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d EntryHash:64ebbb592737e8dcad8bce5db209713a8af18ca306ea6542e86bb1d1c88f5529 Status:TransactionACK} {ChainID:cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d EntryHash:fcfafbef8471d99cfda1cb8d77a15d71d39d96acaf85da836c1f5075a71bb7c3 Status:TransactionACK} {ChainID:6e4540d08d5ac6a1a394e982fb6a2ab8b516ee751c37420055141b94fe070bfe EntryHash:5ecfd88edacbaeea8dde7ead5634dbfb000c37e9d6a07ffa14432f3a1d07549d Status:TransactionACK} {ChainID: EntryHash:d72101ca339aa7206ae6d82594b61051ff901f78276c389ce8b885f34caa9d34 Status:TransactionACK} {ChainID:a642a8674f46696cc47fdb6b65f9c87b2a19c5ea8123b3d2f0c13b6f33a9d5ef EntryHash:2eb17336caf62f1eb869a5380298ced3b2653b496bd90e99a0a2da58c8b3101b Status:TransactionACK}]"

	if received != expected {
		t.Errorf("expected:%s\nreceived:%s", expected, received)
	}
}
