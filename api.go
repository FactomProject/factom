// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	server = "localhost:8088"
)

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

// TODO ...
// ........
