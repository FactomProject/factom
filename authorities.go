// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

type Authority struct {
	AuthorityChainID  string   `json:"identity_chainid"`
	ManagementChainID string   `json:"management_chaind"`
	MatryoshkaHash    string   `json:"matryoshka_hash"`
	SigningKey        string   `json:"signing_key"`
	Status            int      `json:"status"`
	Efficiency        int      `json:"efficiency"`
	CoinbaseAddress   string   `json:"coinbase_address"`
	AnchorKeys        []string `json:"anchor_keys"`
	// TODO: should keyhistory be part of the api return for an Authority?
	// KeyHistory []string `json:"-"`
}

func (a *Authority) String() string {
	var s string

	s += fmt.Sprintln("AuthorityChainID:", a.AuthorityChainID)
	s += fmt.Sprintln("ManagementChainID:", a.ManagementChainID)
	s += fmt.Sprintln("MatryoshkaHash:", a.MatryoshkaHash)
	s += fmt.Sprintln("SigningKey:", a.SigningKey)
	s += fmt.Sprintln("Status:", a.Status)
	s += fmt.Sprintln("Efficiency:", a.Efficiency)
	s += fmt.Sprintln("CoinbaseAddress:", a.CoinbaseAddress)

	s += fmt.Sprintln("AnchorKeys {")
	for _, k := range a.AnchorKeys {
		s += fmt.Sprintln(" ", k)
	}
	s += fmt.Sprintln("}")

	// s += fmt.Sprintln("KeyHisory {")
	// for _, k := range a.KeyHistory {
	// 	s += fmt.Sprintln(" ", k)
	// }
	// s += fmt.Sprintln("}")

	return s
}

// GetAuthorites retrieves a list of the known athorities from factomd
func GetAuthorites() ([]*Authority, error) {
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
