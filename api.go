// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	ed "github.com/agl/ed25519"
)

var (
	server = "localhost:8088"
)

// BuyTestCredits buys Entry Credits for an Entry Credit Key
func BuyTestCredits(key string, amt int) error {
	api := fmt.Sprintf("http://%s/v1/buytestcredits/", server)
	data := url.Values{
		"to":     {key},
		"amount": {strconv.Itoa(amt)},
	}
	resp, err := http.PostForm(api, data)
	if err != nil {
		return err
	}
	resp.Body.Close()
	
	return nil
}

/* TODO finish CommitChain
// CommitChain sends the signed ChainID, the Entry Hash, and the Entry Credit
// public key to the factom network. Once the payment is verified and the
// network is commited to publishing the Chain it may be published by revealing
// the First Entry in the Chain.
func CommitChain(c *Chain, key *[64]byte) error {
	buf := new(bytes.Buffer)
	
	// 1 byte version
	buf.Write([]byte{0})
	
	// 6 byte milliTimestamp (truncated unix time)
	m := milliTime()
	buf.Write(m)

	// 32 byte ChainID Hash
	if c, err := hex.DecodeString(e.ChainID); err != nil {
		return err
	} else {
		// double sha256 hash of ChainID
		h1 := sha256.Sum256(c)
		h2 := sha256.Sum256(h1[:])
		buf.Write(h2[:])
	}
	
	// 32 byte Hash of the Entry Hash + ChainID
	
	// 32 byte Entry Hash of the First Entry
	
	// 1 byte number of Entry Credits to pay
	
	// 32 byte Pubkey
	
	// 64 byte Signature of data from the Verstion to the Entry Credits
}
*/

// CommitEntry sends the signed Entry Hash and the Entry Credit public key to
// the factom network. Once the payment is verified and the network is commited
// to publishing the Entry it may be published with a call to RevealEntry.
func CommitEntry(e *Entry, key *[64]byte) error {
	buf := new(bytes.Buffer)
	
	// 1 byte version
	buf.Write([]byte{0})
	
	// 6 byte milliTimestamp (truncated unix time)
	buf.Write(milliTime())
			
	// 32 byte Entry Hash
	h := e.Hash()
	buf.Write(h[:])
	
	// 1 byte number of entry credits to pay
	if c, err := ecCost(e); err != nil {
		return err
	} else {
		buf.Write([]byte{byte(c)})
	}
	
	// msg is the byte string before the pubkey and sig
	msg := buf.Bytes()
	
	// 32 byte public key
	buf.Write(key[32:64])
	
	// 64 byte signature
	buf.Write(ed.Sign(key, msg)[:])
	
	fmt.Printf("%x\n", buf)
	
	api := fmt.Sprintf("http://%s/v1/commitentry/", server)
	resp, err := http.Post(api, "binary", buf)
	if err != nil {
		return err
	}
	resp.Body.Close()
	
	return nil
}

// GetBlockHeight reports the current Directory Block Height
func GetBlockHeight() (int, error) {
	api := fmt.Sprintf("http://%s/v1/dblockheight/", server)

	resp, err := http.Get(api)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	height, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, err
	}
	return height, nil
}

// GetChain gets a Entry Block Chain by the ChainID. The Chain should contain a
// series of Entry Block Hashes.
func GetChain(hash string) (*Chain, error) {
	chain := new(Chain)
	api := fmt.Sprintf("http://%s/v1/chain/%s", server, hash)

	resp, err := http.Get(api)
	if err != nil {
		return chain, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chain, err
	}

	err = json.Unmarshal(p, chain)
	if err != nil {
		return chain, err
	}

	return chain, nil
}

// GetChains gets all of the Chains. Each Chain should contain a series of
// Entry Block Hashes
func GetChains() ([]Chain, error) {
	chains := make([]Chain, 0)
	api := fmt.Sprintf("http://%s/v1/chains/", server)

	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var chain Chain
		if err := dec.Decode(&chain); err == io.EOF {
			break
		} else if err != nil {
			return chains, err
		}
		chains = append(chains, chain)
	}

	return chains, nil
}

// GetDBlock gets a Directory Block by the Directory Block Hash. The Directory
// Block should contain a series of Entry Block Hashes.
func GetDBlock(hash string) (*DBlock, error) {
	dblock := new(DBlock)
	api := fmt.Sprintf("http://%s/v1/dblock/%s", server, hash)

	resp, err := http.Get(api)
	if err != nil {
		return dblock, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dblock, err
	}

	err = json.Unmarshal(p, dblock)
	if err != nil {
		return dblock, err
	}

	return dblock, nil
}

// GetDBlocks gets the Directory Blocks whithin the Block Height Range provided
// (inclusive). Each DBlock should contain a series of Entry Block Merkel Roots.
func GetDBlocks(from, to int) ([]DBlock, error) {
	dblocks := make([]DBlock, 0)
	api := fmt.Sprintf("http://%s/v1/dblocksbyrange/%s/%s", server,
		strconv.Itoa(from), strconv.Itoa(to))

	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var block DBlock
		if err := dec.Decode(&block); err == io.EOF {
			break
		} else if err != nil {
			return dblocks, err
		}
		dblocks = append(dblocks, block)
	}

	return dblocks, nil
}

// GetDBInfo gets the Directory Block information by the Directory Block Hash.
// The Directory Block Info should contain information about the directory
// block and the BTC transaction containing the Merkle Root.
func GetDBInfo(hash string) (*DBInfo, error) {
	dbinfo := new(DBInfo)
	api := fmt.Sprintf("http://%s/v1/dbinfo/%s", server, hash)

	resp, err := http.Get(api)
	if err != nil {
		return dbinfo, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dbinfo, err
	}
	err = json.Unmarshal(p, dbinfo)
	if err != nil {
		return dbinfo, err
	}

	return dbinfo, nil
}

// GetEBlock gets an entry block specified by the Entry Block Merkel Root. The
// EBlock should contain a series of Entry Hashes.
func GetEBlock(s string) (*EBlock, error) {
	eblock := new(EBlock)
	api := fmt.Sprintf("http://%s/v1/eblock/%s", server, s)

	resp, err := http.Get(api)
	if err != nil {
		return eblock, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eblock, err
	}
	err = json.Unmarshal(p, eblock)
	if err != nil {
		return eblock, err
	}
	return eblock, nil
}

// GetEntry gets an entry based on the Entry Hash. The Entry should contain a
// hex encoded string of Entry Data and a series of External IDs.
func GetEntry(s string) (*Entry, error) {
	entry := new(Entry)
	api := fmt.Sprintf("http://%s/v1/entry/%s", server, s)

	resp, err := http.Get(api)
	if err != nil {
		return entry, err
	}
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entry, err
	}
	err = json.Unmarshal(p, entry)

	return entry, nil
}

func GetEntriesByExtID(eid string) ([]Entry, error) {
	entries := make([]Entry, 0)
	api := fmt.Sprintf("http://%s/v1/entriesbyeid/%s", server, eid)

	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	for {
		var entry Entry
		if err := dec.Decode(&entry); err == io.EOF {
			break
		} else if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func NewECKey() *[64]byte {
	rand, err := os.Open("/dev/random")
	if err != nil {
		return &[64]byte{byte(0)}
	}

	// private key is [32]byte private section + [32]byte public key
	_, priv, err := ed.GenerateKey(rand)
	if err != nil {
		return &[64]byte{byte(0)} 
	}
	return priv
}

func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

func ecCost(e *Entry) (int8, error) {
	p, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	} 
	// n is the capacity of the entry payment in KB
	r := len(p) % 1000
	n := int8(len(p) / 1000)
	if r > 0 {
		n += 1
	}
	if n > 10 {
		return n, fmt.Errorf("Entry larger than 10KB")
	}
	return n, nil
}
