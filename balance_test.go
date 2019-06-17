// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/FactomProject/factom"

	"testing"
)

func TestGetMultipleFCTBalances(t *testing.T) {
	factomdResponse := `{
	  "jsonrpc": "2.0",
	  "id": 3,
	  "result": {
	    "currentheight": 192663,
	    "lastsavedheight": 192662,
	    "balances": [
	      {
	        "ack": 4008,
	        "saved": 4008,
	        "err": ""
	      }, {
	        "ack": 4008,
	        "saved": 4008,
	        "err": ""
	      }, {
	        "ack": 4,
	        "saved": 4,
	        "err": ""
	      }
	    ]
	  }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	fas := []string{
		"FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu",
		"FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu",
		"FA3upjWMKHmStAHR5ZgKVK4zVHPb8U74L2wzKaaSDQEonHajiLeq",
	}
	bs, err := GetMultipleFCTBalances(fas...)
	if err != nil {
		t.Error(err)
	}
	t.Log(bs)
}

func TestGetMultipleECBalances(t *testing.T) {
	factomdResponse := `{
	  "jsonrpc": "2.0",
	  "id": 4,
	  "result": {
	    "currentheight": 192663,
	    "lastsavedheight": 192662,
	    "balances": [
	      {
	        "ack": 4008,
	        "saved": 4008,
	        "err": ""
	      }, {
	        "ack": 4008,
	        "saved": 4008,
	        "err": ""
	      }, {
	        "ack": 4,
	        "saved": 4,
	        "err": ""
	      }
	    ]
	  }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	ecs := []string{
		"EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa",
		"EC1m9mouvUQeEidmqpUYpYtXg8fvTYi6GNHaKg8KMLbdMBrFfmUa",
		"EC3htx3MxKqKTrTMYj4ApWD8T3nYBCQw99veRvH1FLFdjgN6GuNK",
	}
	bs, err := GetMultipleECBalances(ecs...)
	if err != nil {
		t.Error(err)
	}
	t.Log(bs)
}

func TestGetECBalance(t *testing.T) {
	factomdResponse := `{
      "jsonrpc": "2.0",
      "id": 0,
      "result": {
        "balance": 2000
      }
    }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

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
	factomdResponse := `{
      "jsonrpc": "2.0",
      "id": 0,
      "result": {
        "balance": 966582271
      }
    }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, _ := GetFactoidBalance("FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q")

	//fmt.Println(response)
	expectedResponse := int64(966582271)

	if expectedResponse != response {
		fmt.Println(response)
		fmt.Println(expectedResponse)
		t.Fail()
	}
}
