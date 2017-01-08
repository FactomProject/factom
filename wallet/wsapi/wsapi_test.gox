// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// run ``factom-walletd -w /tmp/test_wallet-01 -p 8889`` and ``factomd`` to test the wsapi calls.

package wsapi_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

var testnet = "localhost:8889"

func TestAllAddresses(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":0,"method":"all-addresses"}`

	resp, err := apiCall(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGenerateECAddress(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":0,"method":"generate-ec-address"}`

	resp, err := apiCall(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGenerateFCTAddress(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":0,"method":"generate-factoid-address"}`

	resp, err := apiCall(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestImportAddresses(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":0,"method":"import-addresses","params":{"addresses":[{"secret":"Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"}]}}`

	resp, err := apiCall(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestAddress(t *testing.T) {
	req := `{"jsonrpc":"2.0","id":0,"method":"address","params":{"address":"FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"}}`

	resp, err := apiCall(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestTransaction(t *testing.T) {
	new := `{"jsonrpc":"2.0","id":0,"method":"new-transaction","params":{"tx-name":"a"}}`

	if resp, err := apiCall(new); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	addin := `{"jsonrpc":"2.0","id":0,"method":"add-input","params":{"tx-name":"a","address":"FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q","amount":1000000000}}`
	if resp, err := apiCall(addin); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	addout := `{"jsonrpc":"2.0","id":0,"method":"add-output","params":{"tx-name":"a","address":"FA2xWmGckzbACp7LR43bfnViCyN4uNhugBxiaYPtWGVZVSn1y2m6","amount":1000000000}}`
	if resp, err := apiCall(addout); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	addfee := `{"jsonrpc":"2.0","id":0,"method":"add-fee","params":{"tx-name":"a","address":"FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"}}`
	if resp, err := apiCall(addfee); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	sign := `{"jsonrpc":"2.0","id":0,"method":"sign-transaction","params":{"tx-name":"a"}}`
	if resp, err := apiCall(sign); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	compose := `{"jsonrpc":"2.0","id":0,"method":"compose-transaction","params":{"tx-name":"a"}}`
	if resp, err := apiCall(compose); err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func apiCall(req string) (string, error) {
	client := &http.Client{}
	buf := bytes.NewBuffer([]byte(req))
	re, err := http.NewRequest("POST", "http://"+testnet+"/v2", buf)
	if err != nil {
		return "", err
	}
	re.SetBasicAuth(RpcUser, RpcPass)
	resp, err := client.Do(re)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(p), nil
}
