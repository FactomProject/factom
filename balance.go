// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Balance struct {
	Balance int
}

func ECBalance(key string) (*Balance, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/entry-credit-balance/%s", server, key))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	b := new(Balance)
	if err := json.Unmarshal(body, b); err != nil {
		return nil, err
	}
	
	return b, nil
}