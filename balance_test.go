// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"

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
