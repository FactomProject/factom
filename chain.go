// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Chain struct {
	//chainid was originally required as a paramater passed with the json.
	//it is now overwritten with the chainid derived from the extid elements
	ChainID    string `json:"chainid"`
	FirstEntry *Entry `json:"firstentry"`
}

func NewChain(e *Entry) *Chain {
	c := new(Chain)
	c.FirstEntry = e

	// create the chainid from a series of hashes of the Entries ExtIDs
	hs := sha256.New()
	for _, id := range e.ExtIDs {
		h := sha256.Sum256(id)
		hs.Write(h[:])
	}
	c.ChainID = hex.EncodeToString(hs.Sum(nil))
	c.FirstEntry.ChainID = c.ChainID

	return c
}

func ChainExists(chainid string) bool {
	if _, err := GetChainHead(chainid); err == nil {
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
