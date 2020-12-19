// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/FactomProject/factom"

	"testing"
)

func TestUnmarshalABlock(t *testing.T) {
	// original	// js := []byte(`{"ablock":{"header":{"prevbackrefhash":"e3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f5","dbheight":20000,"headerexpansionsize":0,"headerexpansionarea":"","messagecount":2,"bodysize":131,"adminchainid":"000000000000000000000000000000000000000000000000000000000000000a","chainid":"000000000000000000000000000000000000000000000000000000000000000a"},"abentries":[{"adminidtype":1,"identityadminchainid":"0000000000000000000000000000000000000000000000000000000000000000","prevdbsig":{"pub":"0426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613a","sig":"a7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d"}},{"adminidtype":0,"minutenumber":1}],"backreferencehash":"c8ad13a2aea0f961bf73ac9e79ae8aa0d77ddf59e7d02931de7b9e53a3a20c5e","lookuphash":"e7eb4bda495dbe7657cae1525b6be78bd2fdbad952ebde506b6a97e1cf8f431e"},"rawdata":"000000000000000000000000000000000000000000000000000000000000000ae3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f500004e200000000002000000830100000000000000000000000000000000000000000000000000000000000000000426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613aa7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d0001"}`)
	js := []byte(`{"ablock":{"header":{"prevbackrefhash":"e3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f5","dbheight":20000,"headerexpansionsize":0,"headerexpansionarea":"","messagecount":2,"bodysize":131,"adminchainid":"000000000000000000000000000000000000000000000000000000000000000a","chainid":"000000000000000000000000000000000000000000000000000000000000000a"},"abentries":[{"adminidtype":1,"identityadminchainid":"0000000000000000000000000000000000000000000000000000000000000000","prevdbsig":{"pub":"0426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613a","sig":"a7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d"}},{"adminidtype":0,"minutenumber":1},{"adminidtype":2,"identitychainid":"1111111111111111111111111111111111111111111111111111111111111111","mhash":"2222222222222222222222222222222222222222222222222222222222222222"},{"adminidtype":3,"identitychainid":"3333333333333333333333333333333333333333333333333333333333333333","mhash":"4444444444444444444444444444444444444444444444444444444444444444"},{"adminidtype":4,"amount":3},{"adminidtype":5,"identitychainid":"5555555555555555555555555555555555555555555555555555555555555555","dbheight":10},{"adminidtype":6,"identitychainid":"6666666666666666666666666666666666666666666666666666666666666666","dbheight":11},{"adminidtype":7,"identitychainid":"7777777777777777777777777777777777777777777777777777777777777777","dbheight":12},{"adminidtype":8,"identitychainid":"8888888888888888888888888888888888888888888888888888888888888888","keypriority":2,"publickey":"9999999999999999999999999999999999999999999999999999999999999999","dbheight":13},{"adminidtype":9,"identitychainid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","keypriority":3,"keytype":1,"ecdsapublickey":"13hikQwGStt6Urgiwbxecv5NVW6f77LP3N"},{"adminidtype":13,"identitychainid":"1313131313131313131313131313131313131313131313131313131313131313","factoidaddress":"FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu"},{"adminidtype":14,"identitychainid":"1414141414141414141414141414141414141414141414141414141414141414","efficiency":14}],"backreferencehash":"c8ad13a2aea0f961bf73ac9e79ae8aa0d77ddf59e7d02931de7b9e53a3a20c5e","lookuphash":"e7eb4bda495dbe7657cae1525b6be78bd2fdbad952ebde506b6a97e1cf8f431e"},"rawdata":"000000000000000000000000000000000000000000000000000000000000000ae3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f500004e200000000002000000830100000000000000000000000000000000000000000000000000000000000000000426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613aa7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d0001"}`)

	wrap := new(struct {
		ABlock  *ABlock `json:"ablock"`
		RawData string  `json:"rawdata"`
	})
	if err := json.Unmarshal(js, wrap); err != nil {
		t.Error(err)
	}
	t.Log("ABlock:", wrap.ABlock)
	t.Log("RawData:", wrap.RawData)
}

func TestGetABlock(t *testing.T) {
	factomdResponse := `{
	    "jsonrpc": "2.0",
	    "id": 1,
	    "result": {
	        "ablock": {
	            "header": {
	                "prevbackrefhash": "e3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f5",
	                "dbheight": 20000,
	                "headerexpansionsize": 0,
	                "headerexpansionarea": "",
	                "messagecount": 2,
	                "bodysize": 131,
	                "adminchainid": "000000000000000000000000000000000000000000000000000000000000000a",
	                "chainid": "000000000000000000000000000000000000000000000000000000000000000a"
	            },
	            "abentries": [{
	                "adminidtype": 1,
	                "identityadminchainid": "0000000000000000000000000000000000000000000000000000000000000000",
	                "prevdbsig": {
	                    "pub": "0426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613a",
	                    "sig": "a7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d"
	                }
	            }, {
	                "adminidtype":0,
	                "minutenumber": 1
	            }],
	            "backreferencehash": "c8ad13a2aea0f961bf73ac9e79ae8aa0d77ddf59e7d02931de7b9e53a3a20c5e",
	            "lookuphash": "e7eb4bda495dbe7657cae1525b6be78bd2fdbad952ebde506b6a97e1cf8f431e"
	        },
	        "rawdata": "000000000000000000000000000000000000000000000000000000000000000ae3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f500004e200000000002000000830100000000000000000000000000000000000000000000000000000000000000000426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613aa7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d0001"
	    }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, factomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	ab, err := GetABlock("e7eb4bda495dbe7657cae1525b6be78bd2fdbad952ebde506b6a97e1cf8f431e")
	if err != nil {
		t.Error(err)
	}
	t.Log("ABlock:", ab)
}

func TestGetABlockByHeight(t *testing.T) {
	simlatedFactomdResponse := `{
	    "jsonrpc": "2.0",
	    "id": 1,
	    "result": {
	        "ablock": {
	            "header": {
	                "prevbackrefhash": "e3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f5",
	                "dbheight": 20000,
	                "headerexpansionsize": 0,
	                "headerexpansionarea": "",
	                "messagecount": 2,
	                "bodysize": 131,
	                "adminchainid": "000000000000000000000000000000000000000000000000000000000000000a",
	                "chainid": "000000000000000000000000000000000000000000000000000000000000000a"
	            },
	            "abentries": [{
	                "adminidtype": 1,
	                "identityadminchainid": "0000000000000000000000000000000000000000000000000000000000000000",
	                "prevdbsig": {
	                    "pub": "0426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613a",
	                    "sig": "a7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d"
	                }
	            }, {
	                "adminidtype": 0,
	                "minutenumber": 1
	            }],
	            "backreferencehash": "c8ad13a2aea0f961bf73ac9e79ae8aa0d77ddf59e7d02931de7b9e53a3a20c5e",
	            "lookuphash": "e7eb4bda495dbe7657cae1525b6be78bd2fdbad952ebde506b6a97e1cf8f431e"
	        },
	        "rawdata": "000000000000000000000000000000000000000000000000000000000000000ae3549cd600cbb00d6f8bf4c505ee74f6dc5326d7aa02bb7e4b33f8f16bd6f3f500004e200000000002000000830100000000000000000000000000000000000000000000000000000000000000000426a802617848d4d16d87830fc521f4d136bb2d0c352850919c2679f189613aa7d55725393d78a0e623141a41bfcb64956d308eeb1ae501243ad171c2ed42e62a654e138025d0439ecb5bbf594315c191fa88eedb699d9b63a426a6036d630d0001"
	    }
	}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, simlatedFactomdResponse)
	}))
	defer ts.Close()

	SetFactomdServer(ts.URL[7:])

	ab, err := GetABlockByHeight(20000)
	if err != nil {
		t.Error(err)
	}
	t.Log("ABlock:", ab)
}
