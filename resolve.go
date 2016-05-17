// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

type ResolveAddressResponse struct {
	FactoidAddress     string `json:"factoidaddress"`
	EntryCreditAddress string `json:"entrycreditaddress"`
}

func ResolveDnsName(addr string) (fct, ec string, err error) {
	req := NewJSON2Request("resolve-address", apiCounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return "", "", err
	}
	if resp.Error != nil {
		return "", "", resp.Error
	}

	b := new(ResolveAddressResponse)
	if err := json.Unmarshal(resp.JSONResult(), b); err != nil {
		return "", "", err
	}

	return b.FactoidAddress, b.EntryCreditAddress, nil
}
