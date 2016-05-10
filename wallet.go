// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"strings"
)

type GenerateAddressResponse struct {
	Address string
}

type GenerateAddressFromPrivateKeyRequest struct {
	Name       string `json:"name"`
	PrivateKey string `json:"privateKey,omitempty"`
	Mnemonic   string `json:"mnemonic,omitempty"`
}

func GenerateFactoidAddress(name string) (string, error) {
	name = strings.TrimSpace(name)

	req := NewJSON2Request("factoid-generate-address", apiCounter(), name)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateEntryCreditAddress(name string) (string, error) {
	name = strings.TrimSpace(name)

	req := NewJSON2Request("factoid-generate-ec-address", apiCounter(), name)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateFactoidAddressFromPrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)
	param := GenerateAddressFromPrivateKeyRequest{Name: name, PrivateKey: privateKey}

	req := NewJSON2Request("factoid-generate-address-from-private-key", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateEntryCreditAddressFromPrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)
	param := GenerateAddressFromPrivateKeyRequest{Name: name, PrivateKey: privateKey}

	req := NewJSON2Request("factoid-generate-ec-address-from-private-key", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateFactoidAddressFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)
	param := GenerateAddressFromPrivateKeyRequest{Name: name, PrivateKey: privateKey}

	req := NewJSON2Request("factoid-generate-address-from-human-readable-private-key", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateEntryCreditAddressFromHumanReadablePrivateKey(name string, privateKey string) (string, error) {
	name = strings.TrimSpace(name)
	param := GenerateAddressFromPrivateKeyRequest{Name: name, PrivateKey: privateKey}

	req := NewJSON2Request("factoid-generate-ec-address-from-human-readable-private-key", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

func GenerateFactoidAddressFromMnemonic(name string, mnemonic string) (string, error) {
	name = strings.TrimSpace(name)
	param := GenerateAddressFromPrivateKeyRequest{Name: name, Mnemonic: mnemonic}

	req := NewJSON2Request("factoid-generate-address-from-token-sale", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	add := new(GenerateAddressResponse)
	if err := json.Unmarshal(resp.Result, add); err != nil {
		return "", err
	}

	return add.Address, nil
}

/*
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


*/
