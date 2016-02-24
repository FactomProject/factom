// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type EBlock struct {
	Header struct {
		BlockSequenceNumber int
		ChainID             string
		PrevKeyMR           string
		Timestamp           uint64
	}
	EntryList []EBEntry
}

type EBEntry struct {
	Timestamp int64
	EntryHash string
}

func GetEBlock(keymr string) (*EBlock, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/entry-block-by-keymr/%s", server, keymr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	e := new(EBlock)
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *EBlock) String() string {
	var s string
	s += fmt.Sprintln("BlockSequenceNumber:", e.Header.BlockSequenceNumber)
	s += fmt.Sprintln("ChainID:", e.Header.ChainID)
	s += fmt.Sprintln("PrevKeyMR:", e.Header.PrevKeyMR)
	s += fmt.Sprintln("Timestamp:", e.Header.Timestamp)
	for _, v := range e.EntryList {
		s += fmt.Sprintln("EBEntry {")
		s += fmt.Sprintln("	Timestamp", v.Timestamp)
		s += fmt.Sprintln("	EntryHash", v.EntryHash)
		s += fmt.Sprintln("}")
	}
	return s
}
