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

func TestGetTPS(t *testing.T) {
	factomdResponse := `{
	   "jsonrpc": "2.0",
	   "id": 1,
	   "result": {
	      "totaltxrate": 314.1592,
	      "instanttxrate": 271.828
	   }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	instant, total, err := GetTPS()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Instant: %f, Total: %f\n", instant, total)
}
