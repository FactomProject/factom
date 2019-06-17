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

// TestGetCurrentMinute relies on having a running factom daemon to provide an
// api endpoint at localhost:8088
func TestGetCurrentMinute(t *testing.T) {
	factomdResponse := `{
	   "jsonrpc": "2.0",
	   "id": 1,
	   "result": {
	      "leaderheight": 191244,
	      "directoryblockheight": 191244,
	      "minute": 0,
	      "currentblockstarttime": 1557936697742751200,
	      "currentminutestarttime": 1557936697742751200,
	      "currenttime": 1557936697763826700,
	      "directoryblockinseconds": 600,
	      "stalldetected": false,
	      "faulttimeout": 120,
	      "roundtimeout": 30
	   }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	min, err := GetCurrentMinute()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(min.String())
}
