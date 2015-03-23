package factom

import (
	"fmt"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	server = "localhost:8088"
)

// GetBlockHeight reports the current Block Height
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

// GetDBlocks gets the Directory Blocks whithin the Block Height Range provided
// (inclusive). Each DBlock should contain a series of Entry Block Hashes.
func GetDBlocks(from, to int) ([]DBlock, error) {
	dblocks := make([]DBlock, 0)
	api := fmt.Sprintf("http://%s/v1/dblocksbyrange/%s/%s", server, strconv.Itoa(from),
		strconv.Itoa(to))

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

// GetEBlock gets an entry block specified by the Entry Block Hash. The EBlock
// should contain a series of Entry Hashes
func GetEBlock(s string) (EBlock, error) {
	var eblock EBlock
	api := fmt.Sprintf("http://%s/v1/eblockbymr/%s", server, s)

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



// TODO ...
// ........

// GetChain gets a series of Entry Hashes associated with the Chain ID provided
func GetChain(s string) Chain {
	var c Chain
	return c
}

// GetDBlock gets a Directory Block by the Directory Block Hash. The DBlock
// should contain a series of Entry Block Hashes.
func GetDBlock(s string) DBlock {
	var d DBlock
	return d
}

// GetEntry gets an entry based on the Entry Hash. The Entry should contain a
// hex encoded string of Entry Data and a series of External IDs.
func GetEntry(s string) Entry {
	var e Entry
	return e
}
