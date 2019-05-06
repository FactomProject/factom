// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
)

var (
	ErrChainPending = errors.New("Chain not yet included in a Directory Block")
)

// A Chain is a blockchain datastructure in Factom. The Chain is defined by its
// First Entry from wich the ChainID is derived. Every Entry in the Chain will
// share the ChainID and may be found searching the Factom Entry Blocks.
type Chain struct {
	//chainid was originally required as a paramater passed with the json.
	//it is now overwritten with the chainid derived from the extid elements
	ChainID    string `json:"chainid"`
	FirstEntry *Entry `json:"firstentry"`
}

// NewChain creates a new Factom Chain from an Entry.
func NewChain(e *Entry) *Chain {
	c := new(Chain)
	c.FirstEntry = e

	c.ChainID = ChainIDFromFields(e.ExtIDs)
	c.FirstEntry.ChainID = c.ChainID

	return c
}

// NewChainFromBytes creates a new Factom Chain from byte data used to construct an Entry.
func NewChainFromBytes(content []byte, extids ...[]byte) *Chain {
	e := NewEntryFromBytes(nil, content, extids...)
	c := NewChain(e)
	return c
}

// NewChainFromStrings creates a new Factom Chain from strings used to construct an Entry.
func NewChainFromStrings(content string, extids ...string) *Chain {
	e := NewEntryFromStrings("", content, extids...)
	c := NewChain(e)
	return c
}

// ChainExists returns true if a Chain with the given chainid exists within the
// Factom Blockchain.
func ChainExists(chainid string) bool {
	if _, _, err := GetChainHead(chainid); err == nil {
		// no error means we found the Chain
		return true
	}
	return false
}

// ComposeChainCommit creates a JSON2Request to commit a new Chain via the
// factomd web api. The request includes the marshaled MessageRequest with the
// Entry Credit Signature.
func ComposeChainCommit(c *Chain, ec *ECAddress) (*JSON2Request, error) {
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
	if d, err := EntryCost(e); err != nil {
		return nil, err
	} else {
		buf.WriteByte(byte(d + 10))
	}

	// 32 byte Entry Credit Address Public Key + 64 byte Signature
	sig := ec.Sign(buf.Bytes())
	buf.Write(ec.PubBytes())
	buf.Write(sig[:])

	params := messageRequest{Message: hex.EncodeToString(buf.Bytes())}
	req := NewJSON2Request("commit-chain", APICounter(), params)

	return req, nil
}

// ComposeChainReveal creates a JSON2Request to reveal the Chain via the factomd
// web api.
func ComposeChainReveal(c *Chain) (*JSON2Request, error) {
	p, err := c.FirstEntry.MarshalBinary()
	if err != nil {
		return nil, err
	}
	params := entryRequest{Entry: hex.EncodeToString(p)}

	req := NewJSON2Request("reveal-chain", APICounter(), params)
	return req, nil
}

// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, ec *ECAddress) (string, error) {
	type commitResponse struct {
		Message string `json:"message"`
		TxID    string `json:"txid"`
	}

	req, err := ComposeChainCommit(c, ec)
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

// RevealChain sends the Chain data to the factom network to create a chain that
// has previously been commited.
func RevealChain(c *Chain) (string, error) {
	type revealResponse struct {
		Message string `json:"message"`
		Entry   string `json:"entryhash"`
	}

	req, err := ComposeChainReveal(c)
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

// GetChainHead returns the hash of the most recent Entry made into a given
// Factom Chain.
func GetChainHead(chainid string) (string, bool, error) {
	params := chainIDRequest{ChainID: chainid}
	req := NewJSON2Request("chain-head", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", false, err
	}
	if resp.Error != nil {
		return "", false, resp.Error
	}

	head := new(struct {
		ChainHead          string `json:"chainhead"`
		ChainInProcessList bool   `json:"chaininprocesslist"`
	})
	if err := json.Unmarshal(resp.JSONResult(), head); err != nil {
		return "", false, err
	}

	return head.ChainHead, head.ChainInProcessList, nil
}

// GetAllChainEntries returns a list of all Factom Entries for a given Chain.
func GetAllChainEntries(chainid string) ([]*Entry, error) {
	es := make([]*Entry, 0)

	head, inPL, err := GetChainHead(chainid)
	if err != nil {
		return es, err
	}

	if head == "" && inPL {
		return nil, ErrChainPending
	}

	for ebhash := head; ebhash != ZeroHash; {
		eb, err := GetEBlock(ebhash)
		if err != nil {
			return es, err
		}
		s, err := GetAllEBlockEntries(ebhash)
		if err != nil {
			return es, err
		}
		es = append(s, es...)

		ebhash = eb.Header.PrevKeyMR
	}

	return es, nil
}

// GetAllChainEntriesAtHeight returns a list of all Factom Entries for a given
// Chain at a given point in the Chain's history.
func GetAllChainEntriesAtHeight(chainid string, height int64) ([]*Entry, error) {
	es := make([]*Entry, 0)

	head, inPL, err := GetChainHead(chainid)
	if err != nil {
		return es, err
	}

	if head == "" && inPL {
		return nil, ErrChainPending
	}

	for ebhash := head; ebhash != ZeroHash; {
		eb, err := GetEBlock(ebhash)
		if err != nil {
			return es, err
		}
		if eb.Header.DBHeight > height {
			ebhash = eb.Header.PrevKeyMR
			continue
		}
		s, err := GetAllEBlockEntries(ebhash)
		if err != nil {
			return es, err
		}
		es = append(s, es...)

		ebhash = eb.Header.PrevKeyMR
	}

	return es, nil
}

// GetFirstEntry returns the first Entry used to create the given Factom Chain.
func GetFirstEntry(chainid string) (*Entry, error) {
	e := new(Entry)

	head, inPL, err := GetChainHead(chainid)
	if err != nil {
		return e, err
	}

	if head == "" && inPL {
		return nil, ErrChainPending
	}

	eb, err := GetEBlock(head)
	if err != nil {
		return e, err
	}

	for eb.Header.PrevKeyMR != ZeroHash {
		ebhash := eb.Header.PrevKeyMR
		eb, err = GetEBlock(ebhash)
		if err != nil {
			return e, err
		}
	}

	return GetEntry(eb.EntryList[0].EntryHash)
}
