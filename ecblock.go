// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
		BodyHash       string `json:"bodyhash"`
		PrevHeaderHash string `json:"prevheaderhash"`
		PrevFullHash   string `json:"prevfullhash"`
		DBHeight       int64  `json:"dbheight"`
	} `json:"header"`
	HeaderHash string     `json:"headerhash"`
	FullHash   string     `json:"fullhash"`
	Entries    []ECBEntry `json:"body"`
}

func (e *ECBlock) String() string {
	var s string

	s += fmt.Sprintln("HeaderHash:", e.HeaderHash)
	s += fmt.Sprintln("PrevHeaderHash:", e.Header.PrevHeaderHash)
	s += fmt.Sprintln("FullHash:", e.FullHash)
	s += fmt.Sprintln("PrevFullHash:", e.Header.PrevFullHash)
	s += fmt.Sprintln("BodyHash:", e.Header.BodyHash)
	s += fmt.Sprintln("DBHeight:", e.Header.DBHeight)

	s += fmt.Sprintln("Entries:")
	for _, v := range e.Entries {
		s += fmt.Sprintln(v)
	}

	return s
}

func (e *ECBlock) UnmarshalJSON(js []byte) error {
	tmp := new(struct {
		Header struct {
			BodyHash       string `json:"bodyhash"`
			PrevHeaderHash string `json:"prevheaderhash"`
			PrevFullHash   string `json:"prevfullhash"`
			DBHeight       int64  `json:"dbheight"`
		} `json:"header"`
		HeaderHash string `json:"headerhash"`
		FullHash   string `json:"fullhash"`
		Body       struct {
			Entries []json.RawMessage `json:"entries"`
		} `json:"body"`
	})

	err := json.Unmarshal(js, tmp)
	if err != nil {
		return err
	}

	e.Header.BodyHash = tmp.Header.BodyHash
	e.Header.PrevHeaderHash = tmp.Header.PrevHeaderHash
	e.Header.PrevFullHash = tmp.Header.PrevFullHash
	e.Header.DBHeight = tmp.Header.DBHeight
	e.HeaderHash = tmp.HeaderHash
	e.FullHash = tmp.FullHash
	for _, v := range tmp.Body.Entries {
		switch {
		case regexp.MustCompile(`"number":`).MatchString(string(v)):
			a := new(MinuteNumber)
			err := json.Unmarshal(v, a)
			if err != nil {
				return err
			}
			e.Entries = append(e.Entries, a)
		case regexp.MustCompile(`"serverindexnumber":`).MatchString(string(v)):
			a := new(ServerIndexNumber)
			err := json.Unmarshal(v, a)
			if err != nil {
				return err
			}
			e.Entries = append(e.Entries, a)
		case regexp.MustCompile(`"entryhash":`).MatchString(string(v)):
			if regexp.MustCompile(`"chainidhash":`).MatchString(string(v)) {
				a := new(ChainCommit)
				err := json.Unmarshal(v, a)
				if err != nil {
					return err
				}
				e.Entries = append(e.Entries, a)

			} else {
				a := new(EntryCommit)
				err := json.Unmarshal(v, a)
				if err != nil {
					return err
				}
				e.Entries = append(e.Entries, a)
			}
		default:
		}
	}

	return nil
}

// Entry Credit Block Entries are individual members of the Entry Credit Block.
type ECBEntry interface {
	Type() ECID
	String() string
}

type ServerIndexNumber struct {
	ServerIndexNumber int `json:"serverindexnumber"`
}

func (i *ServerIndexNumber) Type() ECID {
	return ECIDServerIndexNumber
}

func (i *ServerIndexNumber) String() string {
	return fmt.Sprintln("ServerIndexNumber:", i.ServerIndexNumber)
}

type MinuteNumber struct {
	Number int `json:"number"`
}

func (m *MinuteNumber) Type() ECID {
	return ECIDMinuteNumber
}

func (m *MinuteNumber) String() string {
	return fmt.Sprintln("MinuteNumber:", m.Number)
}

type ChainCommit struct {
	Version     int    `json:"version"`
	MilliTime   int64  `json:"millitime"`
	ChainIDHash string `json:"chainidhash"`
	Weld        string `json:"weld"`
	EntryHash   string `json:"entryhash"`
	Credits     int    `json:"credits"`
	ECPubKey    string `json:"ecpubkey"`
	Sig         string `json:"sig"`
}

// TODO: func (c *ChainCommit) UnmarshalJSON(js []byte) error {

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

type EntryCommit struct {
	Version   int    `json:"version"`
	MilliTime int64  `json:"millitime"`
	EntryHash string `json:"entryhash"`
	Credits   int    `json:"credits"`
	ECPubKey  string `json:"ecpubkey"`
	Sig       string `json:"sig"`
}

func (e *EntryCommit) UnmarshalJSON(js []byte) error {
	tmp := new(struct {
		Version   int    `json:"version"`
		MilliTime string `json:"millitime"`
		EntryHash string `json:"entryhash"`
		Credits   int    `json:"credits"`
		ECPubKey  string `json:"ecpubkey"`
		Sig       string `json:"sig"`
	})

	err := json.Unmarshal(js, tmp)
	if err != nil {
		return err
	}

	m := make([]byte, 8)
	if p, err := hex.DecodeString(tmp.MilliTime); err != nil {
		return err
	} else {
		copy(m, p)
	}
	e.MilliTime = int64(binary.BigEndian.Uint64(m))

	e.Version = tmp.Version
	e.EntryHash = tmp.EntryHash
	e.Credits = tmp.Credits
	e.ECPubKey = tmp.ECPubKey
	e.Sig = tmp.Sig

	return nil
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

	// create a wraper construct for the ECBlock API return
	wrap := new(struct {
		ECBlock *ECBlock `json:"ecblock"`
	})
	if err := json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return nil, err
	}

	return wrap.ECBlock, nil
}
