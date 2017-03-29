// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	//"bytes"
	"encoding/hex"
	//"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/FactomProject/factom"
)

func TestGetECBalance(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "balance": 2000
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetECBalance("EC3MAHiZyfuEb5fZP2fSp2gXMv8WemhQEUFXyQ2f2HjSkYx7xY1S")

	//fmt.Println(response)
	expectedResponse := int64(2000)

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetFactoidBalance(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "balance": 966582271
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetFactoidBalance("FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q")

	//fmt.Println(response)
	expectedResponse := int64(966582271)

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetRate(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "rate": 95369
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetRate()

	//fmt.Println(response)
	expectedResponse := uint64(95369)

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetDBlock(t *testing.T) {
	simlatedFactomdResponse := `{  
   "jsonrpc":"2.0",
   "id":0,
   "result":{  
      "header":{  
         "prevblockkeymr":"7d15d82e70201e960655ce3e7cf475c9da593dfb82c6dca6377349bd148bf001",
         "sequencenumber":72497,
         "timestamp":1484858820
      },
      "entryblocklist":[  
         {  
            "chainid":"000000000000000000000000000000000000000000000000000000000000000a",
            "keymr":"3faa880a97ef6ce1feca643cffa015dd6be6a597b3f9260e408c5ac9351d1f8d"
         },
         {  
            "chainid":"000000000000000000000000000000000000000000000000000000000000000c",
            "keymr":"5f8c98930a1874a46b47b65b9376a02fbff65b760f6866519799d69e2bc019ee"
         },
         {  
            "chainid":"000000000000000000000000000000000000000000000000000000000000000f",
            "keymr":"8c6fed0f41317cc45201b5b170a9ac5bc045029e39a90b6061211be2c0678718"
         }
      ]
   }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetDBlock("36c360817761e0d92af464f7c2e94a7495104d6b0a6051218cc53e52d3d519b6")

	//fmt.Println(response)
	expectedResponse := `PrevBlockKeyMR: 7d15d82e70201e960655ce3e7cf475c9da593dfb82c6dca6377349bd148bf001
Timestamp: 1484858820
SequenceNumber: 72497
EntryBlock {
	ChainID 000000000000000000000000000000000000000000000000000000000000000a
	KeyMR 3faa880a97ef6ce1feca643cffa015dd6be6a597b3f9260e408c5ac9351d1f8d
}
EntryBlock {
	ChainID 000000000000000000000000000000000000000000000000000000000000000c
	KeyMR 5f8c98930a1874a46b47b65b9376a02fbff65b760f6866519799d69e2bc019ee
}
EntryBlock {
	ChainID 000000000000000000000000000000000000000000000000000000000000000f
	KeyMR 8c6fed0f41317cc45201b5b170a9ac5bc045029e39a90b6061211be2c0678718
}
`

	if expectedResponse != response.String() {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetDBlockHead(t *testing.T) {
	simlatedFactomdResponse := `{  
   "jsonrpc":"2.0",
   "id":0,
   "result":{  
      "keymr":"7ed5d5b240973676c4a8a71c08c0cedb9e0ea335eaef22995911bcdc0fe9b26b"
   }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetDBlockHead()

	//fmt.Println(response)
	expectedResponse := `7ed5d5b240973676c4a8a71c08c0cedb9e0ea335eaef22995911bcdc0fe9b26b`

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetHeights(t *testing.T) {
	simlatedFactomdResponse := `{  
   "jsonrpc":"2.0",
   "id":0,
   "result":{  
      "directoryblockheight":72498,
      "leaderheight":72498,
      "entryblockheight":72498,
      "entryheight":72498,
      "missingentrycount":0,
      "entryblockdbheightprocessing":72498,
      "entryblockdbheightcomplete":72498
   }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetHeights()

	//fmt.Println(response)
	expectedResponse := `DirectoryBlockHeight: 72498
LeaderHeight: 72498
EntryBlockHeight: 72498
EntryHeight: 72498
`

	if expectedResponse != response.String() {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetEntry(t *testing.T) {
	simlatedFactomdResponse := `{  
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
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetEntry("be5216cc7a5a3ad44b49245aec298f47cbdfca9862dee13b0093e5880012b771")

	//fmt.Println(response)
	expectedResponse := `EntryHash: 1c840bc18be182e89e12f9e63fb8897d13b071b631ced7e656837ccea8fdb3ae
ChainID: df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
ExtID: FactomAnchorChain
Content:
hello world
`

	if expectedResponse != response.String() {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetEBlock(t *testing.T) {
	simlatedFactomdResponse := `{"jsonrpc":"2.0","id":0,"result":{"header":{"blocksequencenumber":35990,"chainid":"df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604","prevkeymr":"7bd1725aa29c988f8f3486512a01976807a0884d4c71ac08d18d1982d905a27a","timestamp":1487042760,"dbheight":75893},"entrylist":[{"entryhash":"cefd9554e9d89132a327e292649031e7b6ccea1cebd80d8a4722e56d0147dd58","timestamp":1487043240},{"entryhash":"61a7f9256f330e50ddf92b296c00fa679588854affc13c380e9945b05fc8e708","timestamp":1487043240}]}}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, _ := GetEBlock("5117490532e46037f8eb660c4fd49cae2a734fc9096b431b2a9a738d7d278398")

	expectedResponse := `BlockSequenceNumber: 35990
ChainID: df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604
PrevKeyMR: 7bd1725aa29c988f8f3486512a01976807a0884d4c71ac08d18d1982d905a27a
Timestamp: 1487042760
DBHeight: 75893
EBEntry {
	Timestamp 1487043240
	EntryHash cefd9554e9d89132a327e292649031e7b6ccea1cebd80d8a4722e56d0147dd58
}
EBEntry {
	Timestamp 1487043240
	EntryHash 61a7f9256f330e50ddf92b296c00fa679588854affc13c380e9945b05fc8e708
}
`

	if expectedResponse != response.String() {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

func TestGetRaw(t *testing.T) {
	simlatedFactomdResponse := `{"jsonrpc":"2.0","id":0,"result":{"data":"df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604181735e2bc1caa844d66bd8ffd4b67e879d22f5b92c1a823008a8266b6bf4954eacdbae3b324a32cd77849bf5ab95782e5d9d8dfcba7c2b627da0d927ae19f3bee16802b7455d628a68c12b3513b75ccf0e67c6e722345fcfa2466f320e5762800008c950001130600000003e47fe17ea16474444d3895d6048b2ade4c71114f9742d31a6e1d7d035019e2ee51d3a04c2e8e4d86b84a22ac3f3a6e90046c28373b34678831fa7c460b7c69570000000000000000000000000000000000000000000000000000000000000002"}}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	p, err := GetRaw("7bd1725aa29c988f8f3486512a01976807a0884d4c71ac08d18d1982d905a27a")
	if err != nil {
		t.Error(err)
	}
	response := hex.EncodeToString(p)

	expectedResponse := `df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604181735e2bc1caa844d66bd8ffd4b67e879d22f5b92c1a823008a8266b6bf4954eacdbae3b324a32cd77849bf5ab95782e5d9d8dfcba7c2b627da0d927ae19f3bee16802b7455d628a68c12b3513b75ccf0e67c6e722345fcfa2466f320e5762800008c950001130600000003e47fe17ea16474444d3895d6048b2ade4c71114f9742d31a6e1d7d035019e2ee51d3a04c2e8e4d86b84a22ac3f3a6e90046c28373b34678831fa7c460b7c69570000000000000000000000000000000000000000000000000000000000000002`

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}

/*

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
*/
