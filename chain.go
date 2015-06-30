// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Chain struct {
	ChainID    string
	FirstEntry *Entry
}

func NewChain(e *Entry) *Chain {
	c := new(Chain)
	c.FirstEntry = e

	// create the chainid from a series of hashes of the Entries ExtIDs
	hs := sha256.New()
	for _, id := range e.ExtIDs {
		p, _ := hex.DecodeString(id)
		h := sha256.Sum256(p)
		hs.Write(h[:])
	}
	c.ChainID = hex.EncodeToString(hs.Sum(nil))
	c.FirstEntry.ChainID = c.ChainID
	
	return c
}

type ChainHead struct {
	EntryBlockKeyMR string
}

// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, name string) error {
	type walletcommit struct {
		Message string
	}
	
	buf := new(bytes.Buffer)

	// 1 byte version
	buf.Write([]byte{0})

	// 6 byte milliTimestamp
	buf.Write(milliTime())

	e := c.FirstEntry

	// 32 byte ChainID Hash
	if p, err := hex.DecodeString(c.ChainID); err != nil {
		return err
	} else {
		// double sha256 hash of ChainID
		buf.Write(shad(p))
	}

	// 32 byte Weld; sha256(sha256(EntryHash + ChainID))
	if cid, err := hex.DecodeString(c.ChainID); err != nil {
		return err
	} else {
		s := append(e.Hash(), cid...)
		buf.Write(shad(s))
	}

	// 32 byte Entry Hash of the First Entry
	buf.Write(e.Hash())

	// 1 byte number of Entry Credits to pay
	if d, err := entryCost(e); err != nil {
		return err
	} else {
		buf.WriteByte(byte(d + 10))
	}

	com := new(walletcommit)
	com.Message = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return err
	}
	fmt.Println("DEBUG: Sending unsigned chain commit to wallet:", string(j))
	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/commit-chain/%s", serverFct, name),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	resp.Body.Close()
    
	return nil
}

func RevealChain(c *Chain) error {
	type reveal struct {
		Entry string
	}

	r := new(reveal)
	if p, err := c.FirstEntry.MarshalBinary(); err != nil {
		return err
	} else {
		r.Entry = hex.EncodeToString(p)
	}

	j, err := json.Marshal(r)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/reveal-chain/", server),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func GetChainHead(chainid string) (*ChainHead, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/chain-head/%s", server, chainid))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	c := new(ChainHead)
	if err := json.Unmarshal(body, c); err != nil {
		return nil, err
	}
	
	return c, nil
}
