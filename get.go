// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

func GetDBlock(keymr string) (*DBlock, error) {
	req := NewJSON2Request("directory-block-by-keymr", apiCounter(), keymr)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	
	dblock := resp.Result.(*DBlock)

	return dblock, nil
}

//func GetDBlockHead() (string, error) {
//	req := NewJSON2Request("directory-block-head", apiCounter(), "")
//	resp, err := factomdRequest(req)
//	if err != nil {
//		return "", err
//	}
//	if resp.Error != nil {
//		return "", resp.Error
//	}
//
//	return resp.Result.KeyMR, nil
//}

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
	if p, err := json.Marshal(resp.Result); err != nil {
		return nil, err
	} else {
		if err := e.UnmarshalJSON(p); err != nil {
			return nil, err
		}
	}
	return e, nil
}

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

func GetEBlock(keymr string) (*EBlock, error) {
	req := NewJSON2Request("entry-block-by-keymr", apiCounter(), keymr)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result.(*EBlock), nil
}
