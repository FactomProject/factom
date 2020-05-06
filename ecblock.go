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
	ErrECIDUndefined   = errors.New("ECID type undefined")
	ErrUnknownECBEntry = errors.New("Unknown Entry Credit Block Entry type")
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
		BodyHash            string `json:"bodyhash"`
		PrevHeaderHash      string `json:"prevheaderhash"`
		PrevFullHash        string `json:"prevfullhash"`
		DBHeight            int64  `json:"dbheight"`
		HeaderExpansionArea []byte `json:"headerexpansionarea,omitempty"`
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
	if e.Header.HeaderExpansionArea != nil {
		s += fmt.Sprintf("HeaderExpansionArea: %x\n", e.Header.HeaderExpansionArea)
	}

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

	// the entry block entry type is not specified in the json data, so detect
	// the entry type by regex and umarshal into the correct type.
	for _, v := range tmp.Body.Entries {
		switch {
		case regexp.MustCompile(`"serverindexnumber":`).MatchString(string(v)):
			a := new(ECServerIndexNumber)
			err := json.Unmarshal(v, a)
			if err != nil {
				return err
			}
			e.Entries = append(e.Entries, a)
		case regexp.MustCompile(`"number":`).MatchString(string(v)):
			a := new(ECMinuteNumber)
			err := json.Unmarshal(v, a)
			if err != nil {
				return err
			}
			e.Entries = append(e.Entries, a)
		case regexp.MustCompile(`"entryhash":`).MatchString(string(v)):
			if regexp.MustCompile(`"chainidhash":`).MatchString(string(v)) {
				a := new(ECChainCommit)
				err := json.Unmarshal(v, a)
				if err != nil {
					return err
				}
				e.Entries = append(e.Entries, a)
			} else {
				a := new(ECEntryCommit)
				err := json.Unmarshal(v, a)
				if err != nil {
					return err
				}
				e.Entries = append(e.Entries, a)
			}
		case regexp.MustCompile(`"numec":`).MatchString(string(v)):
			a := new(ECBalanceIncrease)
			err := json.Unmarshal(v, a)
			if err != nil {
				return err
			}
			e.Entries = append(e.Entries, a)
		default:
			return ErrUnknownECBEntry
		}
	}

	return nil
}

// an ECBEntry is an individual member of the Entry Credit Block.
type ECBEntry interface {
	Type() ECID
	String() string
}

// ECServerIndexNumber shows the index of the server that acknowledged the
// following ECBEntries.
type ECServerIndexNumber struct {
	ServerIndexNumber int `json:"serverindexnumber"`
}

func (i *ECServerIndexNumber) Type() ECID {
	return ECIDServerIndexNumber
}

func (i *ECServerIndexNumber) String() string {
	return fmt.Sprintln("ServerIndexNumber:", i.ServerIndexNumber)
}

// ECMinuteNumber represents the end of a minute minute [1-10] in the order of
// the ECBEntries.
type ECMinuteNumber struct {
	Number int `json:"number"`
}

func (m *ECMinuteNumber) Type() ECID {
	return ECIDMinuteNumber
}

func (m *ECMinuteNumber) String() string {
	return fmt.Sprintln("MinuteNumber:", m.Number)
}

// ECChainCommit pays for and reserves a new chain in Factom.
type ECChainCommit struct {
	Version     int    `json:"version"`
	MilliTime   int64  `json:"millitime"`
	ChainIDHash string `json:"chainidhash"`
	Weld        string `json:"weld"`
	EntryHash   string `json:"entryhash"`
	Credits     int    `json:"credits"`
	ECPubKey    string `json:"ecpubkey"`
	Sig         string `json:"sig"`
}

func (c *ECChainCommit) UnmarshalJSON(js []byte) error {
	tmp := new(struct {
		Version     int    `json:"version"`
		MilliTime   string `json:"millitime"`
		ChainIDHash string `json:"chainidhash"`
		Weld        string `json:"weld"`
		EntryHash   string `json:"entryhash"`
		Credits     int    `json:"credits"`
		ECPubKey    string `json:"ecpubkey"`
		Sig         string `json:"sig"`
	})

	err := json.Unmarshal(js, tmp)
	if err != nil {
		return err
	}

	// convert 6 byte MilliTime into int64
	m := make([]byte, 8)
	if p, err := hex.DecodeString(tmp.MilliTime); err != nil {
		return err
	} else {
		// copy p into the last 6 bytes
		copy(m[2:], p)
	}
	c.MilliTime = int64(binary.BigEndian.Uint64(m))

	c.Version = tmp.Version
	c.ChainIDHash = tmp.ChainIDHash
	c.Weld = tmp.Weld
	c.EntryHash = tmp.EntryHash
	c.Credits = tmp.Credits
	c.ECPubKey = tmp.ECPubKey
	c.Sig = tmp.Sig

	return nil
}

func (c *ECChainCommit) Type() ECID {
	return ECIDChainCommit
}

func (c *ECChainCommit) String() string {
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

// ECEntryCommit pays for and reserves a new entry in Factom.
type ECEntryCommit struct {
	Version   int    `json:"version"`
	MilliTime int64  `json:"millitime"`
	EntryHash string `json:"entryhash"`
	Credits   int    `json:"credits"`
	ECPubKey  string `json:"ecpubkey"`
	Sig       string `json:"sig"`
}

func (e *ECEntryCommit) UnmarshalJSON(js []byte) error {
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

	// convert 6 byte MilliTime into int64
	m := make([]byte, 8)
	if p, err := hex.DecodeString(tmp.MilliTime); err != nil {
		return err
	} else {
		// copy p into the last 6 bytes
		copy(m[2:], p)
	}
	e.MilliTime = int64(binary.BigEndian.Uint64(m))

	e.Version = tmp.Version
	e.EntryHash = tmp.EntryHash
	e.Credits = tmp.Credits
	e.ECPubKey = tmp.ECPubKey
	e.Sig = tmp.Sig

	return nil
}

func (e *ECEntryCommit) Type() ECID {
	return ECIDEntryCommit
}

func (e *ECEntryCommit) String() string {
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

// ECBalanceIncrease pays for and reserves a new entry in Factom.
type ECBalanceIncrease struct {
	ECPubKey string `json:"ecpubkey"`
	TXID     string `json:"txid"`
	Index    uint64 `json:"index"`
	NumEC    uint64 `json:"numec"`
}

func (e *ECBalanceIncrease) Type() ECID {
	return ECIDBalanceIncrease
}

func (e *ECBalanceIncrease) String() string {
	var s string
	s += fmt.Sprintln("BalanceIncrease {")
	s += fmt.Sprintln("	ECPubKey:", e.ECPubKey)
	s += fmt.Sprintln("	TXID:", e.TXID)
	s += fmt.Sprintln("	Index:", e.TXID)
	s += fmt.Sprintln("	NumEC:", e.NumEC)
	s += fmt.Sprintln("}")
	return s
}

func getECBlock(keymr string, noraw bool) (ecblock *ECBlock, raw []byte, err error) {
	params := keyMRRequest{KeyMR: keymr, NoRaw: noraw}
	req := NewJSON2Request("entrycredit-block", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	// create a wraper construct for the ECBlock API return
	wrap := new(struct {
		ECBlock *ECBlock `json:"ecblock"`
		RawData string   `json:"rawdata"`
	})
	err = json.Unmarshal(resp.JSONResult(), wrap)
	if err != nil {
		return
	}

	raw, err = hex.DecodeString(wrap.RawData)
	if err != nil {
		return
	}

	return wrap.ECBlock, raw, nil
}

// GetECBlock requests a specified Entry Credit Block from the factomd API with the raw data
func GetECBlock(keymr string) (ecblock *ECBlock, raw []byte, err error) {
	return getECBlock(keymr, false)
}

// GetSimpleECBlock requests a specified Entry Credit Block from the factomd API without the raw data
func GetSimpleECBlock(keymr string) (ecblock *ECBlock, err error) {
	ecblock, _, err = getECBlock(keymr, true)
	return
}

func getECBlockByHeight(height int64, noraw bool) (ecblock *ECBlock, raw []byte, err error) {
	params := heightRequest{Height: height, NoRaw: noraw}
	req := NewJSON2Request("ecblock-by-height", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	wrap := new(struct {
		ECBlock *ECBlock `json:"ecblock"`
		RawData string   `json:"rawdata"`
	})
	if err = json.Unmarshal(resp.JSONResult(), wrap); err != nil {
		return
	}

	raw, err = hex.DecodeString(wrap.RawData)
	if err != nil {
		return
	}

	return wrap.ECBlock, raw, nil
}

// GetECBlockByHeight request an Entry Credit Block of a given height from the
// factomd API with the raw data
func GetECBlockByHeight(height int64) (ecblock *ECBlock, raw []byte, err error) {
	return getECBlockByHeight(height, false)
}

// GetECBlockByHeight request an Entry Credit Block of a given height from the
// factomd API without the raw data
func GetSimpleECBlockByHeight(height int64) (ecblock *ECBlock, err error) {
	ecblock, _, err = getECBlockByHeight(height, true)
	return
}
