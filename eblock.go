// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"fmt"

	"github.com/FactomProject/factomd/wsapi"
)

func GetAllEBlockEntries(ebhash string) ([]*Entry, error) {
	es := make([]*Entry, 0)

	eb, err := GetEBlock(ebhash)
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

func GetEBlock(keymr string) (*wsapi.EntryBlockResponse, error) {
	resp, err := CallV2("entry-block-by-keymr", false, keymr)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	return resp.Result.(*wsapi.EntryBlockResponse), nil
}
