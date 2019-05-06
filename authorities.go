// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

type Authority struct {
	AuthorityChainID  string             `json:"chainid"`
	ManagementChainID string             `json:"manageid"`
	MatryoshkaHash    string             `json:"matroyshka"` // [sic]
	SigningKey        string             `json:"signingkey"`
	Status            string             `json:"status"`
	AnchorKeys        []AnchorSigningKey `json:"anchorkeys"`
}

func (a *Authority) String() string {
	var s string

	s += fmt.Sprintln("AuthorityChainID:", a.AuthorityChainID)
	s += fmt.Sprintln("ManagementChainID:", a.ManagementChainID)
	s += fmt.Sprintln("MatryoshkaHash:", a.MatryoshkaHash)
	s += fmt.Sprintln("SigningKey:", a.SigningKey)
	s += fmt.Sprintln("Status:", a.Status)

	s += fmt.Sprintln("AnchorKeys {")
	for _, k := range a.AnchorKeys {
		s += fmt.Sprintln(k)
	}
	s += fmt.Sprintln("}")

	return s
}

type AnchorSigningKey struct {
	BlockChain string `json:"blockchain"`
	KeyLevel   byte   `json:"level"`
	KeyType    byte   `json:"keytype"`
	SigningKey string `json:"key"` //if bytes, it is hex
}

func (k *AnchorSigningKey) String() string {
	var s string

	s += fmt.Sprintln("BlockChain:", k.BlockChain)
	s += fmt.Sprintln("KeyLevel:", k.KeyLevel)
	s += fmt.Sprintln("KeyType:", k.KeyType)
	s += fmt.Sprintln("SigningKey:", k.SigningKey)

	return s
}

// GetAuthorites retrieves a list of the known athorities from factomd
func GetAuthorities() ([]*Authority, error) {
	req := NewJSON2Request("authorities", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	// create a temporary type to unmarshal the json object
	a := new(struct {
		Authorities []*Authority `json:"authorities"`
	})

	if err := json.Unmarshal(resp.JSONResult(), a); err != nil {
		return nil, err
	}

	return a.Authorities, nil
}
