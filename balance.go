// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

type MultiBalanceResponse struct {
	CurrentHeight   int `json:"currentheight"`
	LastSavedHeight int `json:"lastsavedheight"`
	Balances        []struct {
		Ack   int    `json:"ack"`
		Saved int    `json:"saved"`
		Err   string `json:"err"`
	} `json:"balances"`
}

// GetECBalance returns the Entry Credit balance of a given Entry
// Credit Public Address.
func GetECBalance(addr string) (int64, error) {
	type balanceResponse struct {
		Balance int64 `json:"balance"`
	}

	params := addressRequest{Address: addr}
	req := NewJSON2Request("entry-credit-balance", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(balanceResponse)
	if err := json.Unmarshal(resp.JSONResult(), balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}

// GetFactoidBalance returns the balance in factoshi (factoid * 1e8) of a given
// Factoid Public Address.
func GetFactoidBalance(addr string) (int64, error) {
	type balanceResponse struct {
		Balance int64 `json:"balance"`
	}

	params := addressRequest{Address: addr}
	req := NewJSON2Request("factoid-balance", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(balanceResponse)
	if err := json.Unmarshal(resp.JSONResult(), balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}

// GetBalanceTotals return the total value of Factoids and Entry Credits in the
// wallet according to the the server acknowledgement and the value saved in the
// blockchain.
func GetBalanceTotals() (fs, fa, es, ea int64, err error) {
	type multiBalanceResponse struct {
		FactoidAccountBalances struct {
			Ack   int64 `json:"ack"`
			Saved int64 `json:"saved"`
		} `json:"fctaccountbalances"`
		EntryCreditAccountBalances struct {
			Ack   int64 `json:"ack"`
			Saved int64 `json:"saved"`
		} `json:"ecaccountbalances"`
	}

	req := NewJSON2Request("wallet-balances", APICounter(), nil)
	resp, err := walletRequest(req)
	if err != nil {
		return
	}

	balances := new(multiBalanceResponse)
	err = json.Unmarshal(resp.JSONResult(), balances)
	if err != nil {
		return
	}

	fs = balances.FactoidAccountBalances.Saved
	fa = balances.FactoidAccountBalances.Ack
	es = balances.EntryCreditAccountBalances.Saved
	ea = balances.EntryCreditAccountBalances.Ack

	return
}

// GetMultipleFCTBalances returns balances for multiple Factoid Addresses from
// the factomd API.
func GetMultipleFCTBalances(fas ...string) (*MultiBalanceResponse, error) {
	type multiAddressRequest struct {
		Addresses []string `json:"addresses"`
	}

	params := multiAddressRequest{fas}
	req := NewJSON2Request("multiple-fct-balances", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}

	balances := new(MultiBalanceResponse)
	err = json.Unmarshal(resp.JSONResult(), balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

// GetMultipleECBalances returns balances for multiple Entry Credit Addresses
// from the factomd API.
func GetMultipleECBalances(ecs ...string) (*MultiBalanceResponse, error) {
	type multiAddressRequest struct {
		Addresses []string `json:"addresses"`
	}

	params := multiAddressRequest{ecs}
	req := NewJSON2Request("multiple-ec-balances", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}

	balances := new(MultiBalanceResponse)
	err = json.Unmarshal(resp.JSONResult(), balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
