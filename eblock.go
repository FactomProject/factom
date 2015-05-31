// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ed "github.com/agl/ed25519"
)

type Chain struct {
	ChainID    string
	FirstEntry *Entry
}

type ChainHead struct {
	EntryBlockKeyMR string
}

type EBlock struct {
	Header struct {
		BlockSequenceNumber int
		ChainID             string
		PrevKeyMR           string
		TimeStamp           uint64
	}
	EBEntries []EBEntry
}

type EBEntry struct {
	TimeStamp int64
	EntryHash string
}

// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, key *[64]byte) error {
	type commit struct {
		CommitChainMsg string
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

	msg := buf.Bytes()

	// 32 byte Pubkey
	buf.Write(key[32:64])

	// 64 byte Signature of data from the Verstion to the Entry Credits
	buf.Write(ed.Sign(key, msg)[:])

	com := new(commit)
	com.CommitChainMsg = hex.EncodeToString(buf.Bytes())
	j, err := json.Marshal(com)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v1/commit-chain/", server),
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

func GetEBlock(keymr string) (*EBlock, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/entry-block-by-keymr/%s", server, keymr))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	e := new(EBlock)
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}
	
	return e, nil
}