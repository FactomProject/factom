// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

// GetDBlock requests a Directory Block from factomd by its Key Merkel Root
func GetDBlock(keymr string) (*DBlock, error) {
	req := NewJSON2Request("directory-block-by-keymr", apiCounter(), keymr)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	db := new(DBlock)
	if err := json.Unmarshal(resp.Result, db); err != nil {
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
	if err := json.Unmarshal(resp.Result, head); err != nil {
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
	if err := json.Unmarshal(resp.Result, height); err != nil {
		return 0, err
	}

	return int(height.Height), nil
}

// GetEntry requests an Entry from factomd by its Entry Hash
func GetEntry(hash string) (*Entry, error) {
	req := NewJSON2Request("entry-by-hash", apiCounter(), hash)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	e := new(Entry)
	if err := json.Unmarshal(resp.Result, e); err != nil {
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

	head := new(CHead)
	if err := json.Unmarshal(resp.Result, head); err != nil {
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
	req := NewJSON2Request("entry-block-by-keymr", apiCounter(), keymr)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(EBlock)
	if err := json.Unmarshal(resp.Result, eb); err != nil {
		return nil, err
	}

	return eb, nil
}

func GetRaw(keymr string) ([]byte, error) {
	req := NewJSON2Request("get-raw-data", apiCounter(), keymr)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	raw := new(RawData)
	if err := json.Unmarshal(resp.Result, raw); err != nil {
		return nil, err
	}

	return raw.GetDataBytes()
}

func GetECBalance(key string) (int64, error) {
	req := NewJSON2Request("entry-credit-balance", apiCounter(), key)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(BalanceResponse)
	if err := json.Unmarshal(resp.Result, balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}

func GetFctBalance(key string) (int64, error) {
	req := NewJSON2Request("factoid-balance", apiCounter(), key)
	resp, err := factomdRequest(req)
	if err != nil {
		return -1, err
	}
	if resp.Error != nil {
		return -1, resp.Error
	}

	balance := new(BalanceResponse)
	if err := json.Unmarshal(resp.Result, balance); err != nil {
		return -1, err
	}

	return balance.Balance, nil
}
