// Copyright 2016 Factom Foundation
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

func TestGetEBlock(t *testing.T) {
	simlatedFactomdResponse := `{"jsonrpc":"2.0","id":0,"result":{"header":{"blocksequencenumber":35990,"chainid":"df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604","prevkeymr":"7bd1725aa29c988f8f3486512a01976807a0884d4c71ac08d18d1982d905a27a","timestamp":1487042760,"dbheight":75893},"entrylist":[{"entryhash":"cefd9554e9d89132a327e292649031e7b6ccea1cebd80d8a4722e56d0147dd58","timestamp":1487043240},{"entryhash":"61a7f9256f330e50ddf92b296c00fa679588854affc13c380e9945b05fc8e708","timestamp":1487043240}]}}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	response, err := GetEBlock("5117490532e46037f8eb660c4fd49cae2a734fc9096b431b2a9a738d7d278398")
	if err != nil {
		t.Error(err)
	}

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
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}
