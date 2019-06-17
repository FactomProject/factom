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

func TestGetECRate(t *testing.T) {
	factomdResponse := `{
        "jsonrpc": "2.0",
        "id": 0,
        "result": {
            "rate": 95369
        }
    }`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetECRate()
	if err != nil {
		t.Error(err)
	}

	expectedResponse := uint64(95369)

	if expectedResponse != response {
		t.Errorf("expected:%d\nrecieved:%d", expectedResponse, response)
	}
	t.Log(response)
}
