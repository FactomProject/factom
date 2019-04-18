// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

func GetHeights() (*HeightsResponse, error) {
	req := NewJSON2Request("heights", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	heights := new(HeightsResponse)
	if err := json.Unmarshal(resp.JSONResult(), heights); err != nil {
		return nil, err
	}

	return heights, nil
}

func GetProperties() (string, string, string, string, string, string, string, string) {
	type propertiesResponse struct {
		FactomdVersion       string `json:"factomdversion"`
		FactomdVersionErr    string `json:"factomdversionerr"`
		FactomdAPIVersion    string `json:"factomdapiversion"`
		FactomdAPIVersionErr string `json:"factomdapiversionerr"`
		WalletVersion        string `json:"walletversion"`
		WalletVersionErr     string `json:"walletversionerr"`
		WalletAPIVersion     string `json:"walletapiversion"`
		WalletAPIVersionErr  string `json:"walletapiversionerr"`
	}

	props := new(propertiesResponse)
	wprops := new(propertiesResponse)
	req := NewJSON2Request("properties", APICounter(), nil)
	wreq := NewJSON2Request("properties", APICounter(), nil)

	resp, err := factomdRequest(req)
	if err != nil {
		props.FactomdVersionErr = err.Error()
	} else if resp.Error != nil {
		props.FactomdVersionErr = resp.Error.Error()
	} else if jerr := json.Unmarshal(resp.JSONResult(), props); jerr != nil {
		props.FactomdVersionErr = jerr.Error()
	}

	wresp, werr := walletRequest(wreq)

	if werr != nil {
		wprops.WalletVersionErr = werr.Error()
	} else if wresp.Error != nil {
		wprops.WalletVersionErr = wresp.Error.Error()
	} else if jwerr := json.Unmarshal(wresp.JSONResult(), wprops); jwerr != nil {
		wprops.WalletVersionErr = jwerr.Error()
	}

	return props.FactomdVersion, props.FactomdVersionErr, props.FactomdAPIVersion, props.FactomdAPIVersionErr, wprops.WalletVersion, wprops.WalletVersionErr, wprops.WalletAPIVersion, wprops.WalletAPIVersionErr

}
