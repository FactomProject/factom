// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"

	"fmt"
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

// GetDBlock requests a Directory Block from factomd by its Key Merkle Root
func GetDBlock(keymr string) (*DBlock, error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("directory-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	db := new(DBlock)
	if err := json.Unmarshal(resp.JSONResult(), db); err != nil {
		return nil, err
	}

	return db, nil
}

func GetDBlockHead() (string, error) {
	req := NewJSON2Request("directory-block-head", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	head := new(DBHead)
	if err := json.Unmarshal(resp.JSONResult(), head); err != nil {
		return "", err
	}

	return head.KeyMR, nil
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

// GetChainHead only returns the chainhead part of the response, so you are losing information
// returned by the api. GetChainHeadAndStatus returns the full repsonse.
// TODO: Depreciate this call, or make it return an error when the chainhead == ""
//			When (chainhead == "" && err == nil ) the ChainInProcessList == true, and we could
//			return an error indicating there is no chainhead found, but it will be created in the
//			next block.
func GetChainHead(chainid string) (string, error) {
	ch, err := getChainHead(chainid)
	if err != nil {
		return "", err
	}
	return ch.ChainHead, nil
}

type chainHeadResponse struct {
	ChainHead          string `json:"chainhead"`
	ChainInProcessList bool   `json:"chaininprocesslist"`
}

func GetChainHeadAndStatus(chainid string) (*chainHeadResponse, error) {
	return getChainHead(chainid)
}

func getChainHead(chainid string) (*chainHeadResponse, error) {
	params := chainIDRequest{ChainID: chainid}
	req := NewJSON2Request("chain-head", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	head := new(chainHeadResponse)
	if err := json.Unmarshal(resp.JSONResult(), head); err != nil {
		return nil, err
	}

	return head, nil
}

// GetAllEBlockEntries requests every Entry in a specific Entry Block
func GetAllEBlockEntries(keymr string) ([]*Entry, error) {
	es := make([]*Entry, 0)

	eb, err := GetEBlock(keymr)
	if err != nil {
		return es, err
	}

	for _, v := range eb.EntryList {
		e, err := GetEntry(v.EntryHash)
		if err != nil {
			return es, err
		}
		es = append(es, e)
	}

	return es, nil
}

// GetEBlock requests an Entry Block from factomd by its Key Merkle Root
func GetEBlock(keymr string) (*EBlock, error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("entry-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(EBlock)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
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

func GetAllChainEntries(chainid string) ([]*Entry, error) {
	es := make([]*Entry, 0)

	head, err := GetChainHeadAndStatus(chainid)
	if err != nil {
		return es, err
	}

	if head.ChainHead == "" && head.ChainInProcessList {
		return nil, fmt.Errorf("Chain not yet included in a Directory Block")
	}

	for ebhash := head.ChainHead; ebhash != ZeroHash; {
		eb, err := GetEBlock(ebhash)
		if err != nil {
			return es, err
		}
		s, err := GetAllEBlockEntries(ebhash)
		if err != nil {
			return es, err
		}
		es = append(s, es...)

		ebhash = eb.Header.PrevKeyMR
	}

	return es, nil
}

func GetFirstEntry(chainid string) (*Entry, error) {
	e := new(Entry)

	head, err := GetChainHeadAndStatus(chainid)
	if err != nil {
		return e, err
	}

	if head.ChainHead == "" && head.ChainInProcessList {
		return nil, fmt.Errorf("Chain not yet included in a Directory Block")
	}

	eb, err := GetEBlock(head.ChainHead)
	if err != nil {
		return e, err
	}

	for eb.Header.PrevKeyMR != ZeroHash {
		ebhash := eb.Header.PrevKeyMR
		eb, err = GetEBlock(ebhash)
		if err != nil {
			return e, err
		}
	}

	return GetEntry(eb.EntryList[0].EntryHash)
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
