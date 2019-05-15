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

func TestGetProperties(t *testing.T) {
	factomdResponse := `{
       "jsonrpc": "2.0",
       "id": 1,
       "result": {
          "factomdversion": "BuiltWithoutVersion",
          "factomdapiversion": "2.0"
       }
    }`
	walletdResponse := `{
       "jsonrpc": "2.0",
       "id": 2,
       "result": {
          "walletversion": "BuiltWithoutVersion",
          "walletapiversion": "2.0"
       }
    }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	wts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, walletdResponse)
	}))
	defer wts.Close()

	SetWalletServer(wts.URL[7:])

	props, err := GetProperties()
	if err != nil {
		t.Error(err)
	}
	t.Log(props)
}
