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
	ChainID string
	ExtIDs  [][]byte
	Content []byte
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
		ChainID string
		ExtIDs  []string
		Content string
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
		ChainID   string
		ChainName []string
		ExtIDs    []string
		Content   string
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
