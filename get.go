// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// GetDBlock requests a Directory Block from factomd by its Key Merkel Root
func GetDBlock(keymr string) (*DBlock, error) {
	param := KeyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("directory-block-by-keymr", apiCounter(), param)
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
	req := NewJSON2Request("directory-block-head", apiCounter(), nil)
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
	req := NewJSON2Request("directory-block-height", apiCounter(), nil)
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
	param := HashRequest{Hash: hash}
	req := NewJSON2Request("entry-by-hash", apiCounter(), param)
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
	req := NewJSON2Request("chain-head", apiCounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	head := new(ChainHead)
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
	param := KeyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("entry-block-by-keymr", apiCounter(), param)
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
	param := HashRequest{Hash: keymr}
	req := NewJSON2Request("get-raw-data", apiCounter(), param)
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

func GetECBalance(key string) (int64, error) {
	param := AddressRequest{Address: key}
	req := NewJSON2Request("entry-credit-balance", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(BalanceResponse)
	if err := json.Unmarshal(resp.JSONResult(), balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}

func GetFctBalance(key string) (int64, error) {
	param := AddressRequest{Address: key}
	req := NewJSON2Request("factoid-balance", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(BalanceResponse)
	if err := json.Unmarshal(resp.JSONResult(), balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}

func DnsBalance(addr string) (int64, int64, error) {
	fct, ec, err := ResolveDnsName(addr)
	if err != nil {
		return -1, -1, err
	}

	f, err1 := GetFctBalance(fct)
	e, err2 := GetECBalance(ec)
	if err1 != nil || err2 != nil {
		return f, e, fmt.Errorf("%s\n%s\n", err1, err2)
	}

	return f, e, nil
}

func GetFee() (int64, error) {
	req := NewJSON2Request("factoid-fee", apiCounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	fee := new(FeeResponse)
	if err := json.Unmarshal(resp.JSONResult(), fee); err != nil {
		return -1, err
	}

	return fee.Fee, nil
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
