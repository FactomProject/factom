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

// GetDBlock requests a Directory Block from factomd by its Key Merkel Root
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

func GetDBlockHeight() (int, error) {
	req := NewJSON2Request("directory-block-height", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}

	height := new(DirectoryBlockHeightResponse)
	if err := json.Unmarshal(resp.JSONResult(), height); err != nil {
		return 0, err
	}

	return int(height.Height), nil
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

func GetChainHead(chainid string) (string, error) {
	type chainHeadResponse struct {
		ChainHead string `json:"chainhead"`
	}

	params := chainIDRequest{ChainID: chainid}
	req := NewJSON2Request("chain-head", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	head := new(chainHeadResponse)
	if err := json.Unmarshal(resp.JSONResult(), head); err != nil {
		return "", err
	}

	return head.ChainHead, nil
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

// GetEBlock requests an Entry Block from factomd by its Key Merkel Root
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

	head, err := GetChainHead(chainid)
	if err != nil {
		return es, err
	}

	for ebhash := head; ebhash != ZeroHash; {
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

	head, err := GetChainHead(chainid)
	if err != nil {
		return e, err
	}

	eb, err := GetEBlock(head)
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

func GetProperties() (string, string, string, error) {
	type propertiesResponse struct {
		FactomdVersion string `json:"factomdversion"`
		APIVersion     string `json:"apiversion"`
		WalletVersion  string `json:"walletversion"`
	}

	req := NewJSON2Request("properties", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", "", "", err
	}
	if resp.Error != nil {
		return "", "", "", resp.Error
	}

	props := new(propertiesResponse)
	if err := json.Unmarshal(resp.JSONResult(), props); err != nil {
		return "", "", "", err
	}

	wresp, err := walletRequest(req)
	if err != nil {
		return props.FactomdVersion, props.APIVersion, props.WalletVersion, err
	}
	if wresp.Error != nil {
		return props.FactomdVersion, props.APIVersion, props.WalletVersion,
			wresp.Error
	}

	wprops := new(propertiesResponse)
	if err := json.Unmarshal(wresp.JSONResult(), wprops); err != nil {
		return "", "", "", err
	}

	return props.FactomdVersion, props.APIVersion, wprops.WalletVersion, nil
}
