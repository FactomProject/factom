// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

// GetECBalance returns the balance in factoshi (factoid * 1e8) of a given Entry
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

// GetRate returns the number of factoshis per entry credit
func GetRate() (uint64, error) {
	type rateResponse struct {
		Rate uint64 `json:"rate"`
	}

	req := NewJSON2Request("entry-credit-rate", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	rate := new(rateResponse)
	if err := json.Unmarshal(resp.JSONResult(), rate); err != nil {
		return 0, err
	}

	return rate.Rate, nil
}

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

// GetEntry requests an Entry from factomd by its Entry Hash
func GetEntry(hash string) (*Entry, error) {
	params := hashRequest{Hash: hash}
	req := NewJSON2Request("entry", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	e := new(Entry)
	if err := json.Unmarshal(resp.JSONResult(), e); err != nil {
		return nil, err
	}

	return e, nil
}

func GetRaw(keymr string) ([]byte, error) {
	params := hashRequest{Hash: keymr}
	req := NewJSON2Request("raw-data", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	raw := new(RawData)
	if err := json.Unmarshal(resp.JSONResult(), raw); err != nil {
		return nil, err
	}

	return raw.GetDataBytes()
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

func GetPendingEntries() (string, error) {

	req := NewJSON2Request("pending-entries", APICounter(), nil)
	resp, err := factomdRequest(req)

	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", err
	}

	rBytes := resp.JSONResult()

	return string(rBytes), nil
}

func GetPendingTransactions() (string, error) {

	req := NewJSON2Request("pending-transactions", APICounter(), nil)
	resp, err := factomdRequest(req)

	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", err
	}
	//fmt.Println("factom resp=", resp)
	transList := resp.JSONResult()

	return string(transList), nil
}
