// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func ECBalance(key string) (int64, error) {
	str := fmt.Sprintf("http://%s/v1/entry-credit-balance/%s", serverFct, key)
	resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return 0, fmt.Errorf("Error getting the balance of %s", key)
	}

	if !b.Success {
		return 0, fmt.Errorf(b.Response)
	}

	v, err := strconv.ParseInt(b.Response, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Error getting the balance of %s", key)
	}

	return v, nil
}

func FctBalance(key string) (int64, error) {
	str := fmt.Sprintf("http://%s/v1/factoid-balance/%s", serverFct, key)
	resp, err := http.Get(str)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
		return 0, fmt.Errorf("Error getting the balance of %s", key)
	}

	if !b.Success {
		return 0, fmt.Errorf(b.Response)
	}

	v, err := strconv.ParseInt(b.Response, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Error getting the balance of %s", key)
	}

	return v, nil
}

func DnsBalance(addr string) (int64, int64, error) {
	str := fmt.Sprintf("http://%s/v1/resolve-address/%s", serverFct, addr)
	resp, err := http.Get(str)
	if err != nil {
		return 0, 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	resp.Body.Close()

	type x struct {
		Response string
		Success  bool
	}
	a := new(x)
	if err := json.Unmarshal(body, a); err != nil {
		return 0, 0, fmt.Errorf("Error getting the balance of %s", addr)
	}

	if !a.Success {
		return 0, 0, fmt.Errorf(a.Response)
	}

	type y struct {
		Fct, Ec string
	}
	b := new(y)
	if err := json.Unmarshal([]byte(a.Response), b); err != nil {
		return 0, 0, fmt.Errorf("Error getting the balance of %s", addr)
	}
	
	f, err1 := FctBalance(b.Fct)
	e, err2 := ECBalance(b.Ec)
	if err1 != nil || err2 != nil {
		return f, e, err1
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
