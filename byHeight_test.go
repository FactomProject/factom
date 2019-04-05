// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/FactomProject/factom"
)

func TestFBlockByHeight(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "fblock": {
      "bodymr": "e12db6a1945d513f066cab66c94dc5cca1b8f90997b95a47b46e70b1656f764a",
      "prevkeymr": "98f3a6bcd978080fb612359f131fdb73c6119dea1f45c977a3fa30dfd507ebf5",
      "prevledgerkeymr": "fe7734478375e16f92f9e5513dbc3f2e830dd2bb263801ba816129242ce83cfe",
      "exchrate": 95369,
      "dbheight": 14460,
      "transactions": [
        {
          "millitimestamp": 1480284840637,
          "inputs": [],
          "outputs": [],
          "outecs": [],
          "rcds": [],
          "sigblocks": [],
          "blockheight": 0
        },
        {
          "millitimestamp": 1480284848556,
          "inputs": [
            {
              "amount": 201144428,
              "address": "dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f",
              "useraddress": ""
            }
          ],
          "outputs": [
            {
              "amount": 200000000,
              "address": "031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f",
              "useraddress": ""
            }
          ],
          "outecs": [],
          "rcds": [
            "3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"
          ],
          "sigblocks": [
            {
              "signatures": [
                "68fd6905eeb276739b2541398db3b1b06d73f99a50803bac83eafabc24be656e26278af6fe8070c85e861e21c39a56a5a422dd2d58dd65a7eeff849f6d02de04"
              ]
            }
          ],
          "blockheight": 0
        },
        {
          "millitimestamp": 1480284956754,
          "inputs": [
            {
              "amount": 401144428,
              "address": "dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f",
              "useraddress": ""
            }
          ],
          "outputs": [
            {
              "amount": 400000000,
              "address": "031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f",
              "useraddress": ""
            }
          ],
          "outecs": [],
          "rcds": [
            "3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"
          ],
          "sigblocks": [
            {
              "signatures": [
                "363c20508bddf5a9d4762e2496a861a1f03ec0dc50389b836dec898a3b37c33a6f831edf057f48a961b2d336231a78137e7402a0ca3a1d5c186ce2bb79e44907"
              ]
            }
          ],
          "blockheight": 0
        }
      ],
      "chainid": "000000000000000000000000000000000000000000000000000000000000000f",
      "keymr": "cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497",
      "ledgerkeymr": "886747480a30f833a27a819fe4b92fbd617cda028329fd2e4b87c7721ff65dea"
    },
    "rawdata": "000000000000000000000000000000000000000000000000000000000000000fe12db6a1945d513f066cab66c94dc5cca1b8f90997b95a47b46e70b1656f764a98f3a6bcd978080fb612359f131fdb73c6119dea1f45c977a3fa30dfd507ebf5fe7734478375e16f92f9e5513dbc3f2e830dd2bb263801ba816129242ce83cfe00000000000174890000387c000000000b0000071c020158a7da22bd000000020158a7da41ac010100dff4f06cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8fdfaf8400031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac68fd6905eeb276739b2541398db3b1b06d73f99a50803bac83eafabc24be656e26278af6fe8070c85e861e21c39a56a5a422dd2d58dd65a7eeff849f6d02de04020158a7da70a5010100b09dae6cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8fafd7c200031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac7e503a078b7f5bac25333c54a724530b97d25b8c2a83f82fdde249987b8bf4a4ba3da41643b7723102c3506bae86f6347ae94b72e8ca4c14617d8d3ca57df70500020158a7da9f9f010100818fccb26cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f818f86c600031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971aced5c095795dc3a2fdcf7689718d10f32be9c656049f169900eac51a34b3f1cb2a0d35ba0a9440892219be15927a4702f062e29ac35238ece375bda4dea9a2a0500020158a7dace9101010083ddb1806c031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f83dceb9400dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f013b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29fdb846a8f9abcb8a1eda9556a8d5f8ef959d1649fddb5ba1ccd3a66b6bc919d8ed225b1273e671685812725a10a7bbac95f0ed9f128bcb3547b068fe3137e50500020158a7dafd8201010081bfa3f46cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f81bede8800031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac66080fa6968bd4b104f1f7a2b5f5bd30b05b066bca956ef28ae94a17f1e18fd102de60f7612811707fda09df69802f5fe4438c1b7b750b57fcc4684c9235520100020158a7db2c7b010100dff4f06cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8fdfaf8400031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac40316508538da5dda6d029e6d8b7eca7b9a5074def82a1d78a0f04ac82fe1efbe14f01b01a6778d648ee47c3acf6f9ab2c9f044087703d55f9901403b991780500020158a7db5b75010100dff4f06cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8fdfaf8400031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac4f050366992fe0a8fb8939825e1286f096f8f1e3528d23e737a499ba73ac1c041501c77f8baf780b3cb915f630a07ac2c0cb312bfb8c89b27ab550fef070a30500020158a7db8a6e010100b09dae6cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8fafd7c200031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971aca6d367d754fd45454c04fd090a68ba4b02ad0316362a7fcd16f164f56d335e8ebd83a3f6e9f5dc1cae9a3ad4eb48a675b8c8b984a1989d359f2bcad368c9c60900020158a7dbb96001010083ddb1806c031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f83dceb9400dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f013b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29c59178bafe8aca410070ba26fd2d4eb607360f52a1b74b60f55a66a3c65bfc5c9ce9b03f6dd4880246723d8542c93fc935f988bc77cc3b9f3d616b414005730300020158a7dbe85201010081bfa3f46cdfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f81bede8800031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f013f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac363c20508bddf5a9d4762e2496a861a1f03ec0dc50389b836dec898a3b37c33a6f831edf057f48a961b2d336231a78137e7402a0ca3a1d5c186ce2bb79e449070000"
  }
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	url := ts.URL[7:]
	SetFactomdServer(url)

	var height int64 = 1

	returnVal, _ := GetFBlockByHeight(height)
	//fmt.Println(returnVal)

	expectedString := `FBlock: {"bodymr":"e12db6a1945d513f066cab66c94dc5cca1b8f90997b95a47b46e70b1656f764a","chainid":"000000000000000000000000000000000000000000000000000000000000000f","dbheight":14460,"exchrate":95369,"keymr":"cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497","ledgerkeymr":"886747480a30f833a27a819fe4b92fbd617cda028329fd2e4b87c7721ff65dea","prevkeymr":"98f3a6bcd978080fb612359f131fdb73c6119dea1f45c977a3fa30dfd507ebf5","prevledgerkeymr":"fe7734478375e16f92f9e5513dbc3f2e830dd2bb263801ba816129242ce83cfe","transactions":[{"blockheight":0,"inputs":[],"millitimestamp":1480284840637,"outecs":[],"outputs":[],"rcds":[],"sigblocks":[]},{"blockheight":0,"inputs":[{"address":"dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f","amount":201144428,"useraddress":""}],"millitimestamp":1480284848556,"outecs":[],"outputs":[{"address":"031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f","amount":200000000,"useraddress":""}],"rcds":["3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"],"sigblocks":[{"signatures":["68fd6905eeb276739b2541398db3b1b06d73f99a50803bac83eafabc24be656e26278af6fe8070c85e861e21c39a56a5a422dd2d58dd65a7eeff849f6d02de04"]}]},{"blockheight":0,"inputs":[{"address":"dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f","amount":401144428,"useraddress":""}],"millitimestamp":1480284956754,"outecs":[],"outputs":[{"address":"031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f","amount":400000000,"useraddress":""}],"rcds":["3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"],"sigblocks":[{"signatures":["363c20508bddf5a9d4762e2496a861a1f03ec0dc50389b836dec898a3b37c33a6f831edf057f48a961b2d336231a78137e7402a0ca3a1d5c186ce2bb79e44907"]}]}]}
`
	//might fail b/c json ordering is non-deterministic
	if returnVal.String() != expectedString {
		fmt.Println(returnVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}

	expectedRawString := `FBlock: {
      "bodymr": "e12db6a1945d513f066cab66c94dc5cca1b8f90997b95a47b46e70b1656f764a",
      "prevkeymr": "98f3a6bcd978080fb612359f131fdb73c6119dea1f45c977a3fa30dfd507ebf5",
      "prevledgerkeymr": "fe7734478375e16f92f9e5513dbc3f2e830dd2bb263801ba816129242ce83cfe",
      "exchrate": 95369,
      "dbheight": 14460,
      "transactions": [
        {
          "millitimestamp": 1480284840637,
          "inputs": [],
          "outputs": [],
          "outecs": [],
          "rcds": [],
          "sigblocks": [],
          "blockheight": 0
        },
        {
          "millitimestamp": 1480284848556,
          "inputs": [
            {
              "amount": 201144428,
              "address": "dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f",
              "useraddress": ""
            }
          ],
          "outputs": [
            {
              "amount": 200000000,
              "address": "031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f",
              "useraddress": ""
            }
          ],
          "outecs": [],
          "rcds": [
            "3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"
          ],
          "sigblocks": [
            {
              "signatures": [
                "68fd6905eeb276739b2541398db3b1b06d73f99a50803bac83eafabc24be656e26278af6fe8070c85e861e21c39a56a5a422dd2d58dd65a7eeff849f6d02de04"
              ]
            }
          ],
          "blockheight": 0
        },
        {
          "millitimestamp": 1480284956754,
          "inputs": [
            {
              "amount": 401144428,
              "address": "dfda7feae639018161018676f141c5744397278c9021e1e9d36e89656c7abe8f",
              "useraddress": ""
            }
          ],
          "outputs": [
            {
              "amount": 400000000,
              "address": "031cce24bcc43b596af105167de2c03603c20ada3314a7cfb47befcad4883e6f",
              "useraddress": ""
            }
          ],
          "outecs": [],
          "rcds": [
            "3f8f50d848f1973751c5776e2f34ab9acf42f72da96d74acd64d2935d75971ac"
          ],
          "sigblocks": [
            {
              "signatures": [
                "363c20508bddf5a9d4762e2496a861a1f03ec0dc50389b836dec898a3b37c33a6f831edf057f48a961b2d336231a78137e7402a0ca3a1d5c186ce2bb79e44907"
              ]
            }
          ],
          "blockheight": 0
        }
      ],
      "chainid": "000000000000000000000000000000000000000000000000000000000000000f",
      "keymr": "cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497",
      "ledgerkeymr": "886747480a30f833a27a819fe4b92fbd617cda028329fd2e4b87c7721ff65dea"
    }
`
	returnRawVal, _ := GetBlockByHeightRaw("f", height)
	if returnRawVal.String() != expectedRawString {
		fmt.Println(returnRawVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}
}
