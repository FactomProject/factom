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
		return "", fmt.Errorf("Error attempting to create %s", name)
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
		return "", fmt.Errorf("Error attempting to create %s", name)
	}

	if !b.Success {
		return "", fmt.Errorf(b.Response)
	}

	return b.Response, nil
}
