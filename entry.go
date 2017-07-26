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
)

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
	s += fmt.Sprintf("EntryHash: %x\n", e.Hash())
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

// ComposeEntryCommit creates a JSON2Request to commit a new Entry via the
// factomd web api. The request includes the marshaled MessageRequest with the
// Entry Credit Signature.
func ComposeEntryCommit(e *Entry, ec *ECAddress) (*JSON2Request, error) {
	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp (truncated unix time)
	buf.Write(milliTime())

	// 32 byte Entry Hash
	buf.Write(e.Hash())

	// 1 byte number of entry credits to pay
	if c, err := EntryCost(e); err != nil {
		return nil, err
	} else {
		buf.WriteByte(byte(c))
	}

	// 32 byte Entry Credit Address Public Key + 64 byte Signature
	sig := ec.Sign(buf.Bytes())
	buf.Write(ec.PubBytes())
	buf.Write(sig[:])

	params := messageRequest{Message: hex.EncodeToString(buf.Bytes())}
	req := NewJSON2Request("commit-entry", APICounter(), params)

	return req, nil
}

// ComposeEntryReveal creates a JSON2Request to reveal the Entry via the factomd
// web api.
func ComposeEntryReveal(e *Entry) (*JSON2Request, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return nil, err
	}

	params := entryRequest{Entry: hex.EncodeToString(p)}
	req := NewJSON2Request("reveal-entry", APICounter(), params)

	return req, nil
}

// CommitEntry sends the signed Entry Hash and the Entry Credit public key to
// the factom network. Once the payment is verified and the network is commited
// to publishing the Entry it may be published with a call to RevealEntry.
func CommitEntry(e *Entry, ec *ECAddress) (string, error) {
	type commitResponse struct {
		Message string `json:"message"`
		TxID    string `json:"txid"`
	}

	req, err := ComposeEntryCommit(e, ec)
	if err != nil {
		return "", err
	}

	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", resp.Error
	}
	r := new(commitResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return "", err
	}

	return r.TxID, nil
}

func RevealEntry(e *Entry) (string, error) {
	type revealResponse struct {
		Message string `json:"message"`
		Entry   string `json:"entryhash"`
	}

	req, err := ComposeEntryReveal(e)
	if err != nil {
		return "", err
	}

	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	r := new(revealResponse)
	if err := json.Unmarshal(resp.JSONResult(), r); err != nil {
		return "", err
	}
	return r.Entry, nil
}
