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

func TestGetDBlockByHeight(t *testing.T) {
	factomdResponse := `{
	   "jsonrpc": "2.0",
	   "id": 3,
	   "result": {
	      "dblock": {
	         "dbhash": "ba79704908f6e96a0aeeceeedd8591cf0949bc538cd5df69b1be7ea8095ed778",
	         "keymr": "cde346e7ed87957edfd68c432c984f35596f29c7d23de6f279351cddecd5dc66",
	         "headerhash": null,
	         "header": {
	            "version": 0,
	            "networkid": 4203931042,
	            "bodymr": "d0d3ce18a3522d925d6445fc70a3e050d7586106200100c805e3c434c5f9ea35",
	            "prevkeymr": "e0e26f41120e2dcb65f9bb6fb61fdfa1beee29e33d0d2110b0ebdb9d9cc05f9b",
	            "prevfullhash": "4e60ea451c7f7230e0a7606872b4dadb57859b573e3a201db434504c24ad6089",
	            "timestamp": 24019950,
	            "dbheight": 100,
	            "blockcount": 4,
	            "chainid": "000000000000000000000000000000000000000000000000000000000000000d"
	         },
	         "dbentries": [
	            {
	               "chainid": "000000000000000000000000000000000000000000000000000000000000000a",
	               "keymr": "cc03cb3558b6b1acd24c5439fadee6523dd2811af82affb60f056df3374b39ae"
	            }, {
	               "chainid": "000000000000000000000000000000000000000000000000000000000000000c",
	               "keymr": "ed01afb79fafba436984a48876082f58e52fec1ccc2920d708ef64ad3beccbbd"
	            }, {
	               "chainid": "000000000000000000000000000000000000000000000000000000000000000f",
	               "keymr": "d9a1de8b02f686a9d4232fa7c8420aa0d9538969923c8eee812352c402c4db0d"
	            }, {
	               "chainid": "df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604",
	               "keymr": "acf8ceaaf70311a6e84d8d7f8d349e5c7958c896afa1c3a4edee09c1f5a80752"
	            }
	         ]
	      },
	      "rawdata": "00fa92e5a2d0d3ce18a3522d925d6445fc70a3e050d7586106200100c805e3c434c5f9ea35e0e26f41120e2dcb65f9bb6fb61fdfa1beee29e33d0d2110b0ebdb9d9cc05f9b4e60ea451c7f7230e0a7606872b4dadb57859b573e3a201db434504c24ad6089016e83ee0000006400000004000000000000000000000000000000000000000000000000000000000000000acc03cb3558b6b1acd24c5439fadee6523dd2811af82affb60f056df3374b39ae000000000000000000000000000000000000000000000000000000000000000ced01afb79fafba436984a48876082f58e52fec1ccc2920d708ef64ad3beccbbd000000000000000000000000000000000000000000000000000000000000000fd9a1de8b02f686a9d4232fa7c8420aa0d9538969923c8eee812352c402c4db0ddf3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604acf8ceaaf70311a6e84d8d7f8d349e5c7958c896afa1c3a4edee09c1f5a80752"
	   }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	d, err := GetDBlockByHeight(100)
	if err != nil {
		t.Error(err)
	}
	t.Log("dblock:", d)
}

func TestGetDBlockHead(t *testing.T) {
	factomdResponse := `{
	   "jsonrpc":"2.0",
	   "id":0,
	   "result":{
	      "keymr":"7ed5d5b240973676c4a8a71c08c0cedb9e0ea335eaef22995911bcdc0fe9b26b"
	   }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	response, err := GetDBlockHead()
	if err != nil {
		t.Error(err)
	}

	expectedResponse := `7ed5d5b240973676c4a8a71c08c0cedb9e0ea335eaef22995911bcdc0fe9b26b`

	if expectedResponse != response {
		t.Errorf("expected:%s\nrecieved:%s", expectedResponse, response)
	}
	t.Log(response)
}
