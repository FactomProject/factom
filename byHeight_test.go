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

var ()

func TestDBlockByHeight(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "dblock": {
      "header": {
        "version": 0,
        "networkid": 4203931043,
        "bodymr": "7716df6083612597d4ef18a8076c40676cd8e0df8110825c07942ca5d30073b4",
        "prevkeymr": "fa7e6e2d37b012d71111bc4e649f1cb9d6f0321964717d35636e4637699d8da2",
        "prevfullhash": "469c7fdce467d222363d55ac234901d8acc61fd4045ae26dd54dac17c556ef86",
        "timestamp": 24671414,
        "dbheight": 14460,
        "blockcount": 3,
        "chainid": "000000000000000000000000000000000000000000000000000000000000000d"
      },
      "dbentries": [
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000a",
          "keymr": "574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"
        },
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000c",
          "keymr": "2a10f1678b9736f213ef3ac76e4f8aa910e5fed66733aa30dafdc91245157b3b"
        },
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000f",
          "keymr": "cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497"
        }
      ],
      "dbhash": "aa7d881d23aad83425c3f10996999d31c76b51b14946f7aca204c150e81bc6d6",
      "keymr": "18509c431ee852edbe1029d676217a0d9cb4fcc11ef8e9aef27fd6075167120c"
    },
    "rawdata": "00fa92e5a37716df6083612597d4ef18a8076c40676cd8e0df8110825c07942ca5d30073b4fa7e6e2d37b012d71111bc4e649f1cb9d6f0321964717d35636e4637699d8da2469c7fdce467d222363d55ac234901d8acc61fd4045ae26dd54dac17c556ef86017874b60000387c00000003000000000000000000000000000000000000000000000000000000000000000a574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0000000000000000000000000000000000000000000000000000000000000000c2a10f1678b9736f213ef3ac76e4f8aa910e5fed66733aa30dafdc91245157b3b000000000000000000000000000000000000000000000000000000000000000fcbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497"
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

	returnVal, _ := GetDBlockByHeight(height)
	//fmt.Println(returnVal)

	expectedString := `DBlock: {"dbentries":[{"chainid":"000000000000000000000000000000000000000000000000000000000000000a","keymr":"574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"},{"chainid":"000000000000000000000000000000000000000000000000000000000000000c","keymr":"2a10f1678b9736f213ef3ac76e4f8aa910e5fed66733aa30dafdc91245157b3b"},{"chainid":"000000000000000000000000000000000000000000000000000000000000000f","keymr":"cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497"}],"dbhash":"aa7d881d23aad83425c3f10996999d31c76b51b14946f7aca204c150e81bc6d6","header":{"blockcount":3,"bodymr":"7716df6083612597d4ef18a8076c40676cd8e0df8110825c07942ca5d30073b4","chainid":"000000000000000000000000000000000000000000000000000000000000000d","dbheight":14460,"networkid":4203931043,"prevfullhash":"469c7fdce467d222363d55ac234901d8acc61fd4045ae26dd54dac17c556ef86","prevkeymr":"fa7e6e2d37b012d71111bc4e649f1cb9d6f0321964717d35636e4637699d8da2","timestamp":24671414,"version":0},"keymr":"18509c431ee852edbe1029d676217a0d9cb4fcc11ef8e9aef27fd6075167120c"}
`
	//might fail b/c json ordering is non-deterministic
	if returnVal.String() != expectedString {
		fmt.Println(returnVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}

	expectedRawString := `DBlock: {
      "header": {
        "version": 0,
        "networkid": 4203931043,
        "bodymr": "7716df6083612597d4ef18a8076c40676cd8e0df8110825c07942ca5d30073b4",
        "prevkeymr": "fa7e6e2d37b012d71111bc4e649f1cb9d6f0321964717d35636e4637699d8da2",
        "prevfullhash": "469c7fdce467d222363d55ac234901d8acc61fd4045ae26dd54dac17c556ef86",
        "timestamp": 24671414,
        "dbheight": 14460,
        "blockcount": 3,
        "chainid": "000000000000000000000000000000000000000000000000000000000000000d"
      },
      "dbentries": [
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000a",
          "keymr": "574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"
        },
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000c",
          "keymr": "2a10f1678b9736f213ef3ac76e4f8aa910e5fed66733aa30dafdc91245157b3b"
        },
        {
          "chainid": "000000000000000000000000000000000000000000000000000000000000000f",
          "keymr": "cbadd7e280377ad8360a4b309df9d14f56552582c05100145ca3367e50adc497"
        }
      ],
      "dbhash": "aa7d881d23aad83425c3f10996999d31c76b51b14946f7aca204c150e81bc6d6",
      "keymr": "18509c431ee852edbe1029d676217a0d9cb4fcc11ef8e9aef27fd6075167120c"
    }
`
	returnRawVal, _ := GetBlockByHeightRaw("d", height)
	if returnRawVal.String() != expectedRawString {
		fmt.Println(returnRawVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}
}

func TestABlockByHeight(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "ablock": {
      "header": {
        "prevbackrefhash": "77e4fb398e228ec9710c20988647a01e2259a40ab77e27c005baf7f2deae3415",
        "dbheight": 14460,
        "headerexpansionsize": 0,
        "headerexpansionarea": "",
        "messagecount": 4,
        "bodysize": 516,
        "adminchainid": "000000000000000000000000000000000000000000000000000000000000000a",
        "chainid": "000000000000000000000000000000000000000000000000000000000000000a"
      },
      "abentries": [
        {
          "identityadminchainid": "888888e238492b2d723d81f7122d4304e5405b18bd9c7cb22ca6bcbc1aab8493",
          "prevdbsig": {
            "pub": "0186ad82617edf3565d944aa104590eb6adb338e92ee6fcd750c2ab2b2707e25",
            "sig": "5796cd49835088ea0d6b8e4a75611ebc674fb791d6e9ebc7f6e5bb1a5e86fc25a8a7742e8f60870e2cb8523fd122ef54bb95ac94b3676b81e07c921ed2196508"
          }
        },
        {
          "identityadminchainid": "888888fc37fa418395eeccb95ab0a4c64d528b2aeefa0d1632c8a116a0e4f5b1",
          "prevdbsig": {
            "pub": "c845f47df202a649e2262d3da0e35556aab62e361425ad7d2e7813a215c8f277",
            "sig": "a5c976c4d18814916fc893f7b4dee78120d20e0deab2b04df2e3b67c2ea1123224db28559ca6d022822388a5ce41128bf5a09ccbbd02b1c5b17a4152183a3d06"
          }
        },
        {
          "identityadminchainid": "88888815ac8a1ab6b8f57cee67ba15aad23ab7d8e70ffdca064200738c201f74",
          "prevdbsig": {
            "pub": "f18512813300d8c1d11e78216d0640ddcc35156a20b53d5ced351a7d5ad90010",
            "sig": "1051c165d7ad33e1f764bb96e5e661053da381ebd708c8ac137da2a1b6847eac07e83472d4fa6096768c7904760c821e45b5ebe23a691cc5bad1b61937f9e303"
          }
        },
        {
          "identityadminchainid": "888888271203752870ae5e6fa0cf96f93cf14bd052455ad476ab26de1ad2c077",
          "prevdbsig": {
            "pub": "4f2d34f0417297e2e985e0cc6e4cf3d0814416d09f37af7375517ea236786ed3",
            "sig": "01206ff2963af7df29bb6749a4c29bc1eb65a48bd7b3ec6590c723e11c3a3c5342e72f9b079a58d77a2562c25289d799fadfc5205f1e99c4f1d5c3ce85432906"
          }
        }
      ],
      "backreferencehash": "1786de6a72311dd4b60c6608d60c2b9367642fb1ee6b867b2c9f4c57c87b8cba",
      "lookuphash": "574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"
    },
    "rawdata": "000000000000000000000000000000000000000000000000000000000000000a77e4fb398e228ec9710c20988647a01e2259a40ab77e27c005baf7f2deae34150000387c00000000040000020401888888e238492b2d723d81f7122d4304e5405b18bd9c7cb22ca6bcbc1aab84930186ad82617edf3565d944aa104590eb6adb338e92ee6fcd750c2ab2b2707e255796cd49835088ea0d6b8e4a75611ebc674fb791d6e9ebc7f6e5bb1a5e86fc25a8a7742e8f60870e2cb8523fd122ef54bb95ac94b3676b81e07c921ed219650801888888fc37fa418395eeccb95ab0a4c64d528b2aeefa0d1632c8a116a0e4f5b1c845f47df202a649e2262d3da0e35556aab62e361425ad7d2e7813a215c8f277a5c976c4d18814916fc893f7b4dee78120d20e0deab2b04df2e3b67c2ea1123224db28559ca6d022822388a5ce41128bf5a09ccbbd02b1c5b17a4152183a3d060188888815ac8a1ab6b8f57cee67ba15aad23ab7d8e70ffdca064200738c201f74f18512813300d8c1d11e78216d0640ddcc35156a20b53d5ced351a7d5ad900101051c165d7ad33e1f764bb96e5e661053da381ebd708c8ac137da2a1b6847eac07e83472d4fa6096768c7904760c821e45b5ebe23a691cc5bad1b61937f9e30301888888271203752870ae5e6fa0cf96f93cf14bd052455ad476ab26de1ad2c0774f2d34f0417297e2e985e0cc6e4cf3d0814416d09f37af7375517ea236786ed301206ff2963af7df29bb6749a4c29bc1eb65a48bd7b3ec6590c723e11c3a3c5342e72f9b079a58d77a2562c25289d799fadfc5205f1e99c4f1d5c3ce85432906"
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

	returnVal, _ := GetABlockByHeight(height)
	//fmt.Println(returnVal)

	expectedString := `ABlock: {"abentries":[{"identityadminchainid":"888888e238492b2d723d81f7122d4304e5405b18bd9c7cb22ca6bcbc1aab8493","prevdbsig":{"pub":"0186ad82617edf3565d944aa104590eb6adb338e92ee6fcd750c2ab2b2707e25","sig":"5796cd49835088ea0d6b8e4a75611ebc674fb791d6e9ebc7f6e5bb1a5e86fc25a8a7742e8f60870e2cb8523fd122ef54bb95ac94b3676b81e07c921ed2196508"}},{"identityadminchainid":"888888fc37fa418395eeccb95ab0a4c64d528b2aeefa0d1632c8a116a0e4f5b1","prevdbsig":{"pub":"c845f47df202a649e2262d3da0e35556aab62e361425ad7d2e7813a215c8f277","sig":"a5c976c4d18814916fc893f7b4dee78120d20e0deab2b04df2e3b67c2ea1123224db28559ca6d022822388a5ce41128bf5a09ccbbd02b1c5b17a4152183a3d06"}},{"identityadminchainid":"88888815ac8a1ab6b8f57cee67ba15aad23ab7d8e70ffdca064200738c201f74","prevdbsig":{"pub":"f18512813300d8c1d11e78216d0640ddcc35156a20b53d5ced351a7d5ad90010","sig":"1051c165d7ad33e1f764bb96e5e661053da381ebd708c8ac137da2a1b6847eac07e83472d4fa6096768c7904760c821e45b5ebe23a691cc5bad1b61937f9e303"}},{"identityadminchainid":"888888271203752870ae5e6fa0cf96f93cf14bd052455ad476ab26de1ad2c077","prevdbsig":{"pub":"4f2d34f0417297e2e985e0cc6e4cf3d0814416d09f37af7375517ea236786ed3","sig":"01206ff2963af7df29bb6749a4c29bc1eb65a48bd7b3ec6590c723e11c3a3c5342e72f9b079a58d77a2562c25289d799fadfc5205f1e99c4f1d5c3ce85432906"}}],"backreferencehash":"1786de6a72311dd4b60c6608d60c2b9367642fb1ee6b867b2c9f4c57c87b8cba","header":{"adminchainid":"000000000000000000000000000000000000000000000000000000000000000a","bodysize":516,"chainid":"000000000000000000000000000000000000000000000000000000000000000a","dbheight":14460,"headerexpansionarea":"","headerexpansionsize":0,"messagecount":4,"prevbackrefhash":"77e4fb398e228ec9710c20988647a01e2259a40ab77e27c005baf7f2deae3415"},"lookuphash":"574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"}
`
	//might fail b/c json ordering is non-deterministic
	if returnVal.String() != expectedString {
		fmt.Println(returnVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}

	expectedRawString := `ABlock: {
      "header": {
        "prevbackrefhash": "77e4fb398e228ec9710c20988647a01e2259a40ab77e27c005baf7f2deae3415",
        "dbheight": 14460,
        "headerexpansionsize": 0,
        "headerexpansionarea": "",
        "messagecount": 4,
        "bodysize": 516,
        "adminchainid": "000000000000000000000000000000000000000000000000000000000000000a",
        "chainid": "000000000000000000000000000000000000000000000000000000000000000a"
      },
      "abentries": [
        {
          "identityadminchainid": "888888e238492b2d723d81f7122d4304e5405b18bd9c7cb22ca6bcbc1aab8493",
          "prevdbsig": {
            "pub": "0186ad82617edf3565d944aa104590eb6adb338e92ee6fcd750c2ab2b2707e25",
            "sig": "5796cd49835088ea0d6b8e4a75611ebc674fb791d6e9ebc7f6e5bb1a5e86fc25a8a7742e8f60870e2cb8523fd122ef54bb95ac94b3676b81e07c921ed2196508"
          }
        },
        {
          "identityadminchainid": "888888fc37fa418395eeccb95ab0a4c64d528b2aeefa0d1632c8a116a0e4f5b1",
          "prevdbsig": {
            "pub": "c845f47df202a649e2262d3da0e35556aab62e361425ad7d2e7813a215c8f277",
            "sig": "a5c976c4d18814916fc893f7b4dee78120d20e0deab2b04df2e3b67c2ea1123224db28559ca6d022822388a5ce41128bf5a09ccbbd02b1c5b17a4152183a3d06"
          }
        },
        {
          "identityadminchainid": "88888815ac8a1ab6b8f57cee67ba15aad23ab7d8e70ffdca064200738c201f74",
          "prevdbsig": {
            "pub": "f18512813300d8c1d11e78216d0640ddcc35156a20b53d5ced351a7d5ad90010",
            "sig": "1051c165d7ad33e1f764bb96e5e661053da381ebd708c8ac137da2a1b6847eac07e83472d4fa6096768c7904760c821e45b5ebe23a691cc5bad1b61937f9e303"
          }
        },
        {
          "identityadminchainid": "888888271203752870ae5e6fa0cf96f93cf14bd052455ad476ab26de1ad2c077",
          "prevdbsig": {
            "pub": "4f2d34f0417297e2e985e0cc6e4cf3d0814416d09f37af7375517ea236786ed3",
            "sig": "01206ff2963af7df29bb6749a4c29bc1eb65a48bd7b3ec6590c723e11c3a3c5342e72f9b079a58d77a2562c25289d799fadfc5205f1e99c4f1d5c3ce85432906"
          }
        }
      ],
      "backreferencehash": "1786de6a72311dd4b60c6608d60c2b9367642fb1ee6b867b2c9f4c57c87b8cba",
      "lookuphash": "574e7d6178e04c92879601f0cb84a619f984eb2617ff9e76ee830a9f614cc9a0"
    }
`
	returnRawVal, _ := GetBlockByHeightRaw("a", height)
	if returnRawVal.String() != expectedRawString {
		fmt.Println(returnRawVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}
}

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

func TestECBlockByHeight(t *testing.T) {
	simlatedFactomdResponse := `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "ecblock": {
      "header": {
        "bodyhash": "ef7a85d4bf868e34aff4edce479f6ee412161e1faa3596a112cd5ef75e96f59c",
        "prevheaderhash": "add44ed20133c7b8c9500ab5819d3aee665fffcce7acb6baa098fa8210b43a8b",
        "prevfullhash": "2f1dd9e5f1ab34102f65dea55c1598e2344568d68c0511640b7f436295615746",
        "dbheight": 14460,
        "headerexpansionarea": "",
        "objectcount": 10,
        "bodysize": 20,
        "chainid": "000000000000000000000000000000000000000000000000000000000000000c",
        "ecchainid": "000000000000000000000000000000000000000000000000000000000000000c"
      },
      "body": {
        "entries": [
          {
            "number": 1
          },
          {
            "number": 2
          },
          {
            "number": 3
          },
          {
            "number": 4
          },
          {
            "number": 5
          },
          {
            "number": 6
          },
          {
            "number": 7
          },
          {
            "number": 8
          },
          {
            "number": 9
          },
          {
            "number": 10
          }
        ]
      }
    },
    "rawdata": "000000000000000000000000000000000000000000000000000000000000000cef7a85d4bf868e34aff4edce479f6ee412161e1faa3596a112cd5ef75e96f59cadd44ed20133c7b8c9500ab5819d3aee665fffcce7acb6baa098fa8210b43a8b2f1dd9e5f1ab34102f65dea55c1598e2344568d68c0511640b7f4362956157460000387c00000000000000000a0000000000000014010101020103010401050106010701080109010a"
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

	returnVal, _ := GetECBlockByHeight(height)
	//fmt.Println(returnVal)

	expectedString := `ECBlock: {"body":{"entries":[{"number":1},{"number":2},{"number":3},{"number":4},{"number":5},{"number":6},{"number":7},{"number":8},{"number":9},{"number":10}]},"header":{"bodyhash":"ef7a85d4bf868e34aff4edce479f6ee412161e1faa3596a112cd5ef75e96f59c","bodysize":20,"chainid":"000000000000000000000000000000000000000000000000000000000000000c","dbheight":14460,"ecchainid":"000000000000000000000000000000000000000000000000000000000000000c","headerexpansionarea":"","objectcount":10,"prevfullhash":"2f1dd9e5f1ab34102f65dea55c1598e2344568d68c0511640b7f436295615746","prevheaderhash":"add44ed20133c7b8c9500ab5819d3aee665fffcce7acb6baa098fa8210b43a8b"}}
`
	//might fail b/c json ordering is non-deterministic
	if returnVal.String() != expectedString {
		fmt.Println(returnVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}

	expectedRawString := `ECBlock: {
      "header": {
        "bodyhash": "ef7a85d4bf868e34aff4edce479f6ee412161e1faa3596a112cd5ef75e96f59c",
        "prevheaderhash": "add44ed20133c7b8c9500ab5819d3aee665fffcce7acb6baa098fa8210b43a8b",
        "prevfullhash": "2f1dd9e5f1ab34102f65dea55c1598e2344568d68c0511640b7f436295615746",
        "dbheight": 14460,
        "headerexpansionarea": "",
        "objectcount": 10,
        "bodysize": 20,
        "chainid": "000000000000000000000000000000000000000000000000000000000000000c",
        "ecchainid": "000000000000000000000000000000000000000000000000000000000000000c"
      },
      "body": {
        "entries": [
          {
            "number": 1
          },
          {
            "number": 2
          },
          {
            "number": 3
          },
          {
            "number": 4
          },
          {
            "number": 5
          },
          {
            "number": 6
          },
          {
            "number": 7
          },
          {
            "number": 8
          },
          {
            "number": 9
          },
          {
            "number": 10
          }
        ]
      }
    }
`
	returnRawVal, _ := GetBlockByHeightRaw("ec", height)
	if returnRawVal.String() != expectedRawString {
		fmt.Println(returnRawVal.String())
		fmt.Println(expectedString)
		t.Fail()
	}
}
