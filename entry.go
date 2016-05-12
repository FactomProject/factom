// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ed "github.com/FactomProject/ed25519"
)

type EntryStrings struct {
	ChainID string   `json:"chainid"`
	ExtIDs  []string `json:"extids"`
	Content string   `json:"content"`
}

func (es *EntryStrings) ToEntry() (*Entry, error) {
	e := new(Entry)
	e.ChainID = es.ChainID

	for _, v := range es.ExtIDs {
		if p, err := hex.DecodeString(v); err != nil {
			return nil, fmt.Errorf("Could not decode ExtID %s: %s", v, err)
		} else {
			e.ExtIDs = append(e.ExtIDs, p)
		}
	}

	p, err := hex.DecodeString(es.Content)
	if err != nil {
		return nil, fmt.Errorf("Could not decode Content %s: %s", es.Content, err)
	}
	e.Content = p

	return e, nil
}

type Entry struct {
	ChainID string   `json:"chainid"`
	ExtIDs  [][]byte `json:"extids"`
	Content []byte   `json:"content"`
}

func (e *Entry) Hash() []byte {
	a, err := e.MarshalBinary()
	if err != nil {
		return make([]byte, 32)
	}
	return sha52(a)
}

func (e *Entry) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	ids, err := e.MarshalExtIDsBinary()
	if err != nil {
		return buf.Bytes(), err
	}

	// Header

	// 1 byte Version
	buf.Write([]byte{0})

	// 32 byte chainid
	if p, err := hex.DecodeString(e.ChainID); err != nil {
		return buf.Bytes(), err
	} else {
		buf.Write(p)
	}

	// 2 byte size of extids
	if err := binary.Write(buf, binary.BigEndian, int16(len(ids))); err != nil {
		return buf.Bytes(), err
	}

	// Body

	// ExtIDs
	buf.Write(ids)

	// Content
	buf.Write(e.Content)

	return buf.Bytes(), nil
}

func (e *Entry) MarshalExtIDsBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, v := range e.ExtIDs {
		// 2 byte length of extid
		binary.Write(buf, binary.BigEndian, int16(len(v)))
		// extid
		buf.Write(v)
	}

	return buf.Bytes(), nil
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	type js struct {
		ChainID string   `json:"chainid"`
		ExtIDs  []string `json:"extids"`
		Content string   `json:"content"`
	}

	j := new(js)

	j.ChainID = e.ChainID

	for _, id := range e.ExtIDs {
		j.ExtIDs = append(j.ExtIDs, hex.EncodeToString(id))
	}

	j.Content = hex.EncodeToString(e.Content)

	return json.Marshal(j)
}

func (e *Entry) String() string {
	var s string
	s += fmt.Sprintln("ChainID:", e.ChainID)
	for _, id := range e.ExtIDs {
		s += fmt.Sprintln("ExtID:", string(id))
	}
	s += fmt.Sprintln("Content:")
	s += fmt.Sprintln(string(e.Content))
	return s
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	type js struct {
		ChainID   string   `json:"chainid"`
		ChainName []string `json:"chainname"`
		ExtIDs    []string `json:"extids"`
		Content   string   `json:"content"`
	}

	j := new(js)
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}

	e.ChainID = j.ChainID

	if e.ChainID == "" {
		n := new(Entry)
		for _, v := range j.ChainName {
			if p, err := hex.DecodeString(v); err != nil {
				return fmt.Errorf("Could not decode ChainName %s: %s", v, err)
			} else {
				n.ExtIDs = append(n.ExtIDs, p)
			}
		}
		m := NewChain(n)
		e.ChainID = m.ChainID
	}

	for _, v := range j.ExtIDs {
		if p, err := hex.DecodeString(v); err != nil {
			return fmt.Errorf("Could not decode ExtID %s: %s", v, err)
		} else {
			e.ExtIDs = append(e.ExtIDs, p)
		}
	}

	p, err := hex.DecodeString(j.Content)
	if err != nil {
		return fmt.Errorf("Could not decode Content %s: %s", j.Content, err)
	}
	e.Content = p

	return nil
}

func ComposeEntryCommit(pub *[32]byte, pri *[64]byte, e *Entry) ([]byte, error) {
	type commit struct {
		CommitEntryMsg string
	}

	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp (truncated unix time)
	buf.Write(milliTime())

	// 32 byte Entry Hash
	buf.Write(e.Hash())

	// 1 byte number of entry credits to pay
	if c, err := entryCost(e); err != nil {
		return nil, err
	} else {
		buf.WriteByte(byte(c))
	}

	// sign the commit
	sig := ed.Sign(pri, buf.Bytes())

	// 32 byte Entry Credit Public Key
	buf.Write(pub[:])

	// 64 byte Signature
	buf.Write(sig[:])

	com := new(commit)
	com.CommitEntryMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func ComposeEntryReveal(e *Entry) ([]byte, error) {
	type reveal struct {
		Entry string `json:"entry"`
	}

	r := new(reveal)
	if p, err := e.MarshalBinary(); err != nil {
		return nil, err
	} else {
		r.Entry = hex.EncodeToString(p)
	}

	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func ComposeChainCommit(pub *[32]byte, pri *[64]byte, c *Chain) ([]byte, error) {
	type commit struct {
		CommitChainMsg string `json:"commitchainmsg"`
	}

	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp
	buf.Write(milliTime())

	e := c.FirstEntry

	// 32 byte ChainID Hash
	if p, err := hex.DecodeString(c.ChainID); err != nil {
		return nil, err
	} else {
		// double sha256 hash of ChainID
		buf.Write(shad(p))
	}

	// 32 byte Weld; sha256(sha256(EntryHash + ChainID))
	if cid, err := hex.DecodeString(c.ChainID); err != nil {
		return nil, err
	} else {
		s := append(e.Hash(), cid...)
		buf.Write(shad(s))
	}

	// 32 byte Entry Hash of the First Entry
	buf.Write(e.Hash())

	// 1 byte number of Entry Credits to pay
	if d, err := entryCost(e); err != nil {
		return nil, err
	} else {
		buf.WriteByte(byte(d + 10))
	}

	// sign the commit
	sig := ed.Sign(pri, buf.Bytes())

	// 32 byte pubkey
	buf.Write(pub[:])

	// 64 byte Signature
	buf.Write(sig[:])

	com := new(commit)
	com.CommitChainMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// CommitEntry sends the signed Entry Hash and the Entry Credit public key to
// the factom network. Once the payment is verified and the network is commited
// to publishing the Entry it may be published with a call to RevealEntry.
func CommitEntry(e *Entry, name string) error {
	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp (truncated unix time)
	buf.Write(milliTime())

	// 32 byte Entry Hash
	buf.Write(e.Hash())

	// 1 byte number of entry credits to pay
	if c, err := entryCost(e); err != nil {
		return err
	} else {
		buf.WriteByte(byte(c))
	}

	param := EntryRequest{Entry: hex.EncodeToString(buf.Bytes())}
	req := NewJSON2Request("commit-entry", apiCounter(), param)
	resp, err := walletRequest(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}
	/*
		rer := new(RevealEntryResponse)
		if err := json.Unmarshal(resp.Result, rer); err != nil {
			return nil, err
		}
	*/

	return nil
}

type RevealEntryResponse struct {
	Message string `json:"message"`
	TxID    string `json:"txid"`
}

func RevealEntry(e *Entry) (*RevealEntryResponse, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return nil, err
	}

	param := EntryRequest{Entry: hex.EncodeToString(p)}
	req := NewJSON2Request("reveal-entry", apiCounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	rer := new(RevealEntryResponse)
	if err := json.Unmarshal(resp.Result, rer); err != nil {
		return nil, err
	}

	return rer, nil
}
