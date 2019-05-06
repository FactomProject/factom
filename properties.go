// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

// PropertiesResponse represents properties of the running factomd and factom
// wallet.
type PropertiesResponse struct {
	FactomdVersion       string `json:"factomdversion"`
	FactomdVersionErr    string `json:"factomdversionerr"`
	FactomdAPIVersion    string `json:"factomdapiversion"`
	FactomdAPIVersionErr string `json:"factomdapiversionerr"`
	WalletVersion        string `json:"walletversion"`
	WalletVersionErr     string `json:"walletversionerr"`
	WalletAPIVersion     string `json:"walletapiversion"`
	WalletAPIVersionErr  string `json:"walletapiversionerr"`
}

// GetProperties requests various properties of the factomd and factom wallet
// software and API versions.
func GetProperties() (*PropertiesResponse, error) {
	// get properties from the factom API and the wallet API
	props := new(PropertiesResponse)
	// wprops := new(PropertiesResponse)
	req := NewJSON2Request("properties", APICounter(), nil)
	wreq := NewJSON2Request("properties", APICounter(), nil)

	resp, err := factomdRequest(req)
	if err != nil {
		props.FactomdVersionErr = err.Error()
		return props, err
	} else if resp.Error != nil {
		props.FactomdVersionErr = resp.Error.Error()
	} else if jerr := json.Unmarshal(resp.JSONResult(), props); jerr != nil {
		props.FactomdVersionErr = jerr.Error()
		return props, jerr
	}

	wresp, werr := walletRequest(wreq)
	wprops := new(PropertiesResponse)
	if werr != nil {
		props.WalletVersionErr = werr.Error()
		return props, werr
	} else if wresp.Error != nil {
		props.WalletVersionErr = wresp.Error.Error()
	} else if jwerr := json.Unmarshal(wresp.JSONResult(), wprops); jwerr != nil {
		props.WalletVersionErr = jwerr.Error()
		return props, jwerr
	}
	props.WalletVersion = wprops.WalletVersion
	props.WalletVersionErr = wprops.WalletVersionErr
	props.WalletAPIVersion = wprops.WalletAPIVersion
	props.WalletVersionErr = wprops.WalletAPIVersionErr

	return props, nil
}
