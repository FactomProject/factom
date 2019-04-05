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
