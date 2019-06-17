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

func TestGetReciept(t *testing.T) {
	factomdResponse := `{
       "jsonrpc": "2.0",
       "id": 1,
       "result": {
          "receipt": {
             "entry": {
                "entryhash": "96b2b60a0e026f3aac01e1680b4d4205ec696845b1b18a1ab6340e21835b6cfe"
             },
             "merklebranch": [
                {
                   "left": "5f9457e8ad1eb2d7a6f2b640141035e6a1e4389d81ca6e18aab9705a83d42e48",
                   "right": "96b2b60a0e026f3aac01e1680b4d4205ec696845b1b18a1ab6340e21835b6cfe",
                   "top": "df025dd89485a38f69867ecbb18fc2c8ff549d9287765877b73be67d3a31a174"
                }, {
                   "left": "df025dd89485a38f69867ecbb18fc2c8ff549d9287765877b73be67d3a31a174",
                   "right": "858586f758d0ca09e842eaa4cf04bc0bb8892228123f13d2569ae16365aa7750",
                   "top": "73f9351a088d2228d94b6a928a5b45840d5c050ca3ffd6905dc443f4ca7adf03"
                }, {
                   "left": "c006042d665b94b6baa6105305cf02233b007192c294ce5c4c078c843fbb1ebe",
                   "right": "73f9351a088d2228d94b6a928a5b45840d5c050ca3ffd6905dc443f4ca7adf03",
                   "top": "ef7646f2f9251c9e50e19ab9343c25eb88c241aa49b7ca779c2318b8ccce1f8a"
                }, {
                   "left": "df3ade9eec4b08d5379cc64270c30ea7315d8a8a1a69efe2b98a60ecdd69e604",
                   "right": "ef7646f2f9251c9e50e19ab9343c25eb88c241aa49b7ca779c2318b8ccce1f8a",
                   "top": "0542f612db0d11bb7f0b3c2bd20363239fdab526f43c099c502e0e40995fed36"
                }, {
                   "left": "54d807c29273ef48d06d0f6a65cd6587566812157770a2ef032cd92db72d0c07",
                   "right": "0542f612db0d11bb7f0b3c2bd20363239fdab526f43c099c502e0e40995fed36",
                   "top": "3e7a636e0d95e2568005a9fb60ecd2a3c168a5a6fe71d097ac3567c9348cd0c1"
                }, {
                   "left": "43d96e6490c0d2aeeeb836f226e08531f1c357e7508b42a9856fd94353b2e5f4",
                   "right": "3e7a636e0d95e2568005a9fb60ecd2a3c168a5a6fe71d097ac3567c9348cd0c1",
                   "top": "56dbd1e0fb4bd7d13aa4cf1c2a32fe015e62650dcdc0171dc07a9197ffb4af54"
                }, {
                   "left": "2ee84f9404e8bac1a413a4151a76fa655859ebdfe71e9385c41a057b64d02bb0",
                   "right": "56dbd1e0fb4bd7d13aa4cf1c2a32fe015e62650dcdc0171dc07a9197ffb4af54",
                   "top": "3a38fec82b26ee916891dab3dd7a7e101ab643aff4a641d895137a7a7c9cac55"
                }
             ],
             "entryblockkeymr": "ef7646f2f9251c9e50e19ab9343c25eb88c241aa49b7ca779c2318b8ccce1f8a",
             "directoryblockkeymr": "3a38fec82b26ee916891dab3dd7a7e101ab643aff4a641d895137a7a7c9cac55",
             "bitcointransactionhash": "38464b98ffe44c71f063ff3bedf80db5e8bb6fe3848322a0966e50a60a65cfc5",
             "bitcoinblockhash": "00000000000000000589540fdaacf4f6ba37513aedc1033e68a649ffde0573ad"
          }
       }
    }`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	r, err := GetReceipt("96b2b60a0e026f3aac01e1680b4d4205ec696845b1b18a1ab6340e21835b6cfe")
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}
