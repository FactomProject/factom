// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/FactomProject/factom"
)

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
