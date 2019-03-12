// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrECIDUndefined = errors.New("ECID type undefined")
)

// ECID defines the type of an Entry Credit Block Entry
type ECID byte

// Available ECID types
const (
	ECIDServerIndexNumber ECID = iota // 0
	ECIDMinuteNumber                  // 1
	ECIDChainCommit                   // 2
	ECIDEntryCommit                   // 3
	ECIDBalanceIncrease               // 4
)

func (id ECID) String() string {
	switch id {
	case ECIDServerIndexNumber:
		return "ServerIndexNumber"
	case ECIDMinuteNumber:
		return "MinuteNumber"
	case ECIDChainCommit:
		return "ChainCommit"
	case ECIDEntryCommit:
		return "EntryCommit"
	case ECIDBalanceIncrease:
		return "BalanceIncrease"
	default:
		return "ECIDUndefined"
	}
}

// ECBlock (Entry Credit Block) holds transactions that create Chains and
// Entries, and fund Entry Credit Addresses.
type ECBlock struct {
	Header struct {
		PrevHeaderHash string `json:"prevheaderhash"`
		PrevFullHash   string `json:"prevfullhash"`
		DBHeight       int64  `json:"dbheight"`
	} `json:"header"`
	Entries []ECBEntry `json:"entries"`
}

// Entry Credit Block Entries are individual members of the Entry Credit Block.
type ECBEntry interface {
	Type() ECID
	String() string
	UnmarshalJSON([]byte) error
}

func (e *ECBlock) String() string {
	var s string
	s += fmt.Sprintln("PrevHeaderHash:", e.Header.PrevHeaderHash)
	s += fmt.Sprintln("PrevFullHash:", e.Header.PrevFullHash)
	s += fmt.Sprintln("DBHeight:", e.Header.DBHeight)
	for _, v := range e.Entries {
		s += fmt.Sprintln(v.Type(), " {")
		s += fmt.Sprintln(v)
		s += fmt.Sprintln("}")
	}
	return s
}

type ServerIndexNumber struct {
	ServerIndexNumber uint8 `json:"serverindexnumber"`
}

func (i *ServerIndexNumber) Type() ECID {
	return ECIDServerIndexNumber
}

func (i *ServerIndexNumber) String() string {
	return fmt.Sprintln("ServerIndexNumber:", i.ServerIndexNumber)
}

func (i *ServerIndexNumber) UnmarshalJSON(js []byte) error {
	if err := json.Unmarshal(js, i); err != nil {
		return err
	}
	return nil
}

type MinuteNumber struct {
	Number uint8 `json:"number"`
}

func (m *MinuteNumber) Type() ECID {
	return ECIDMinuteNumber
}

func (m *MinuteNumber) String() string {
	return fmt.Sprintln("MinuteNumber:", m.Number)
}

func (m *MinuteNumber) UnmarshalJSON(js []byte) error {
	if err := json.Unmarshal(js, m); err != nil {
		return err
	}
	return nil
}

type ChainCommit struct {
	Version     uint8  `json:"version"`
	MilliTime   int64  `json:"millitime"`
	ChainIDHash string `json:"chainidhash"`
	Weld        string `json:"weld"`
	EntryHash   string `json:"entryhash"`
	Credits     uint8  `json:"credits"`
	ECPubKey    string `json:"ecpubkey"`
	Sig         string `json:"sig"`
}

func (c *ChainCommit) Type() ECID {
	return ECIDChainCommit
}

func (c *ChainCommit) String() string {
	var s string
	s += fmt.Sprintln("ChainCommit {")
	s += fmt.Sprintln("	Version:", c.Version)
	s += fmt.Sprintln("	Millitime:", c.MilliTime)
	s += fmt.Sprintln("	ChainIDHash:", c.ChainIDHash)
	s += fmt.Sprintln("	Weld:", c.Weld)
	s += fmt.Sprintln("	EntryHash:", c.EntryHash)
	s += fmt.Sprintln("	Credits:", c.Credits)
	s += fmt.Sprintln("	ECPubKey:", c.ECPubKey)
	s += fmt.Sprintln("	Signature:", c.Sig)
	s += fmt.Sprintln("}")
	return s
}

func (c *ChainCommit) UnmarshalJSON(js []byte) error {
	if err := json.Unmarshal(js, c); err != nil {
		return err
	}
	return nil
}

type EntryCommit struct {
	Version   uint8  `json:"version"`
	MilliTime int64  `json:"millitime"`
	EntryHash string `json:"entryhash"`
	Credits   uint8  `json:"credits"`
	ECPubKey  string `json:"ecpubkey"`
	Sig       string `json:"sig"`
}

func (e *EntryCommit) Type() ECID {
	return ECIDEntryCommit
}

func (e *EntryCommit) String() string {
	var s string
	s += fmt.Sprintln("EntryCommit {")
	s += fmt.Sprintln("	Version:", e.Version)
	s += fmt.Sprintln("	Millitime:", e.MilliTime)
	s += fmt.Sprintln("	EntryHash:", e.EntryHash)
	s += fmt.Sprintln("	Credits:", e.Credits)
	s += fmt.Sprintln("	ECPubKey:", e.ECPubKey)
	s += fmt.Sprintln("	Signature:", e.Sig)
	s += fmt.Sprintln("}")
	return s
}

func (e *EntryCommit) UnmarshalJSON(js []byte) error {
	if err := json.Unmarshal(js, e); err != nil {
		return err
	}
	return nil
}

func GetECBlock(keymr string) (*ECBlock, error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("entrycredit-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	eb := new(ECBlock)
	if err := json.Unmarshal(resp.JSONResult(), eb); err != nil {
		return nil, err
	}

	return eb, nil
}
