// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/FactomProject/factomd/wsapi"
)

func ECBalance(key string) (int64, error) {
	resp, err := CallV2("entry-credit-balance", false, key, new(wsapi.EntryCreditBalanceResponse))
	if err != nil {
		return 0, err
	}

	if resp.Error != nil {
		return 0, fmt.Errorf(resp.Error.Message)
	}

	return resp.Result.(*wsapi.EntryCreditBalanceResponse).Balance, nil
}

func FctBalance(key string) (int64, error) {
	resp, err := CallV2("factoid-balance", false, key, new(wsapi.FactoidBalanceResponse))
	if err != nil {
		return 0, err
	}

	if resp.Error != nil {
		return 0, fmt.Errorf(resp.Error.Message)
	}

	return resp.Result.(*wsapi.FactoidBalanceResponse).Balance, nil
}

func DnsBalance(addr string) (int64, int64, error) {
	fct, ec, err := ResolveDnsName(addr)
	if err != nil {
		return 0, 0, err
	}

	f, err1 := FctBalance(fct)
	e, err2 := ECBalance(ec)
	if err1 != nil || err2 != nil {
		return f, e, fmt.Errorf("%s\n%s\n", err1, err2)
	}

	return f, e, nil
}

func GenerateFactoidAddress(name string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-address/%s", serverFct, name)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}

	fmt.Println(body)

	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateEntryCreditAddress(name string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-ec-address/%s", serverFct, name)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}

	fmt.Println(string(body))

	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateFactoidAddressFromPrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-address-from-private-key/?name=%s&privateKey=%s", serverFct, name, privateKey)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateEntryCreditAddressFromPrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-ec-address-from-private-key/?name=%s&privateKey=%s", serverFct, name, privateKey)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateFactoidAddressFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-address-from-human-readable-private-key/?name=%s&privateKey=%s", serverFct, name, privateKey)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateEntryCreditAddressFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-ec-address-from-human-readable-private-key/?name=%s&privateKey=%s", serverFct, name, privateKey)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}

func GenerateFactoidAddressFromMnemonic(name string, mnemonic string) (string, error) {
	name = strings.TrimSpace(name)

	mnemonic = strings.Replace(mnemonic, " ", "%20", -1)

	str := fmt.Sprintf("http://%s/v1/factoid-generate-address-from-token-sale/?name=%s&mnemonic=%s", serverFct, name, mnemonic)

	resp, err := http.Get(str)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return "", fmt.Errorf("Error attempting to create %s - %v", name, err)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}
