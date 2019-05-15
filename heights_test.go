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

func TestGetHeights(t *testing.T) {
	factomdResponse := `{
       "jsonrpc":"2.0",
       "id":0,
       "result":{
          "directoryblockheight":72498,
          "leaderheight":72498,
          "entryblockheight":72498,
          "entryheight":72498
       }
    }`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetHeights()
	if err != nil {
		t.Error(err)
	}

	expectedResponse := `DirectoryBlockHeight: 72498
LeaderHeight: 72498
EntryBlockHeight: 72498
EntryHeight: 72498
`

	if expectedResponse != response.String() {
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}
