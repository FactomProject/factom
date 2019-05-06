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

// Tests reqire a local factomd node to be running and servier the API!

func TestGetDBlock(t *testing.T) {
	d, raw, err := GetDBlock("cde346e7ed87957edfd68c432c984f35596f29c7d23de6f279351cddecd5dc66")
	if err != nil {
		t.Error(err)
	}
	t.Log("dblock:", d)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
}

func TestGetDBlockByHeight(t *testing.T) {
	d, raw, err := GetDBlockByHeight(100)
	if err != nil {
		t.Error(err)
	}
	t.Log("dblock:", d)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
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

	response, err := GetDBlockHead()
	if err != nil {
		t.Error(err)
	}

	//fmt.Println(response)
	expectedResponse := `7ed5d5b240973676c4a8a71c08c0cedb9e0ea335eaef22995911bcdc0fe9b26b`

	if expectedResponse != response {
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}
