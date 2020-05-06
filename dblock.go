// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// DBlock is a Factom Network Directory Block containing the Merkel root of all
// of the Entries and blocks from a 10 minute period in the Factom Network. The
// Directory Block Key Merkel Root is anchored into the Bitcoin and other
// blockchains for added security and immutability.
type DBlock struct {
	DBHash         string `json:"dbhash"`
	KeyMR          string `json:"keymr"`
	HeaderHash     string `json:"headerhash"`
	SequenceNumber int64  `json:"sequencenynumber"`
	Header         struct {
		Version      int    `json:"version"`
		NetworkID    int    `json:"networkid"`
		BodyMR       string `json:"bodymr"`
		PrevKeyMR    string `json:"prevkeymr"`
		PrevFullHash string `json:"prevfullhash"`
		Timestamp    int    `json:"timestamp"` //in minutes
		DBHeight     int    `json:"dbheight"`
		BlockCount   int    `json:"blockcount"`
	} `json:"header"`
	DBEntries []struct {
		ChainID string `json:"chainid"`
		KeyMR   string `json:"keymr"`
	} `json:"dbentries"`
}

func (db *DBlock) String() string {
	var s string

	s += fmt.Sprintln("DBHash:", db.DBHash)
	s += fmt.Sprintln("KeyMR:", db.KeyMR)
	s += fmt.Sprintln("HeaderHash:", db.HeaderHash)
	s += fmt.Sprintln("SequenceNumber:", db.SequenceNumber)
	s += fmt.Sprintln("Version:", db.Header.Version)
	s += fmt.Sprintln("NetworkID:", db.Header.NetworkID)
	s += fmt.Sprintln("BodyMR:", db.Header.BodyMR)
	s += fmt.Sprintln("PrevKeyMR:", db.Header.PrevKeyMR)
	s += fmt.Sprintln("PrevFullHash:", db.Header.PrevFullHash)
	s += fmt.Sprintln("Timestamp:", db.Header.Timestamp)
	s += fmt.Sprintln("DBHeight:", db.Header.DBHeight)
	s += fmt.Sprintln("BlockCount:", db.Header.BlockCount)

	s += fmt.Sprintln("DBEntries {")
	for _, v := range db.DBEntries {
		s += fmt.Sprintln("	ChainID:", v.ChainID)
		s += fmt.Sprintln("	KeyMR:", v.KeyMR)
	}
	s += fmt.Sprintln("}")

	return s
}

// TODO: GetDBlock should use the dblock api call directy instead of
// re-directing to dblock-by-height.
// we either need to change the "directoy-block" API call or add a new call to
// return the propper information (it should match the dblock-by-height call)

func getDBlock(keymr string, noraw bool) (dblock *DBlock, raw []byte, err error) {
	params := keyMRRequest{KeyMR: keymr}
	req := NewJSON2Request("directory-block", APICounter(), params) // doesn't have a "no raw" request
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	db := new(struct {
		DBHash string `json:"dbhash"`
		Header struct {
			PrevBlockKeyMR string `json:"prevblockkeymr"`
			SequenceNumber int64  `json:"sequencenumber"`
			Timestamp      int64  `json:"timestamp"`
		} `json:"header"`
		EntryBlockList []struct {
			ChainID string `json:"chainid"`
			KeyMR   string `json:"keymr"`
		} `json:"entryblocklist"`
	})

	err = json.Unmarshal(resp.JSONResult(), db)
	if err != nil {
		return
	}

	// TODO: we need a better api call for dblock by keymr so that API will
	// retrun the same as dblock-byheight
	return getDBlockByHeight(db.Header.SequenceNumber, noraw)
}

// GetDBlock requests a Directory Block by its Key Merkle Root from the factomd
// API with the raw data
func GetDBlock(keymr string) (dblock *DBlock, raw []byte, err error) {
	return getDBlock(keymr, false)
}

// GetSimpleDBlock requests a Directory Block by its Key Merkle Root from the factomd
// API without the raw data
func GetSimpleDBlock(keymr string) (dblock *DBlock, err error) {
	dblock, _, err = getDBlock(keymr, true)
	return
}

func getDBlockByHeight(height int64, noraw bool) (dblock *DBlock, raw []byte, err error) {
	params := heightRequest{Height: height, NoRaw: noraw}
	req := NewJSON2Request("dblock-by-height", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, nil, resp.Error
	}

	wrap := new(struct {
		DBlock  *DBlock `json:"dblock"`
		RawData string  `json:"rawdata"`
	})

	err = json.Unmarshal(resp.JSONResult(), wrap)
	if err != nil {
		return
	}

	raw, err = hex.DecodeString(wrap.RawData)
	if err != nil {
		return
	}

	wrap.DBlock.SequenceNumber = height
	return wrap.DBlock, raw, nil
}

// GetDBlockByHeight requests a Directory Block by its block height from the factomd
// API.
func GetDBlockByHeight(height int64) (dblock *DBlock, raw []byte, err error) {
	return getDBlockByHeight(height, false)
}

// GetSimpleDBlockByHeight requests a Directory Block by its block height from the factomd
// API.
func GetSimpleDBlockByHeight(height int64) (dblock *DBlock, err error) {
	dblock, _, err = getDBlockByHeight(height, true)
	return
}

// GetDBlockHead requests the most recent Directory Block Key Merkel Root
// created by the Factom Network.
func GetDBlockHead() (string, error) {
	req := NewJSON2Request("directory-block-head", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	head := new(struct {
		KeyMR string `json:"keymr"`
	})
	if err := json.Unmarshal(resp.JSONResult(), head); err != nil {
		return "", err
	}

	return head.KeyMR, nil
}

// ReplayDBlockFromHeight requests DBlock states to be emitted over the LiveFeed API
func ReplayDBlockFromHeight(startheight int64, endheight int64) (*replayResponse, error) {
	params := replayRequest{StartHeight: startheight, EndHeight: endheight}
	req := NewJSON2Request("replay-from-height", APICounter(), params)
	resp, err := factomdRequest(req)

	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	finalResp := new(replayResponse)
	err = json.Unmarshal(resp.Result, finalResp)
	if err != nil {
		return nil, err
	}

	return finalResp, nil
}
