// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

type JStruct struct {
	data []byte
}

func (e *JStruct) MarshalJSON() ([]byte, error) {
	return e.data, nil
}

func (e *JStruct) UnmarshalJSON(b []byte) error {
	e.data = b
	return nil
}

type BlockByHeightRawResponse struct {
	//TODO: implement all of the blocks as proper structures

	DBlock  *JStruct `json:"dblock,omitempty"`
	ABlock  *JStruct `json:"ablock,omitempty"`
	FBlock  *JStruct `json:"fblock,omitempty"`
	ECBlock *JStruct `json:"ecblock,omitempty"`

	RawData string `json:"rawdata,omitempty"`
}

func (f *BlockByHeightRawResponse) String() string {
	var s string
	if f.DBlock != nil {
		j, _ := f.DBlock.MarshalJSON()
		s += fmt.Sprintln("DBlock:", string(j))
	} else if f.ABlock != nil {
		j, _ := f.ABlock.MarshalJSON()
		s += fmt.Sprintln("ABlock:", string(j))
	} else if f.FBlock != nil {
		j, _ := f.FBlock.MarshalJSON()
		s += fmt.Sprintln("FBlock:", string(j))
	} else if f.ECBlock != nil {
		j, _ := f.ECBlock.MarshalJSON()
		s += fmt.Sprintln("ECBlock:", string(j))
	}

	return s
}

// GetBlockByHeightRaw fetches the specified block type by height
// Deprecated: use ablock, dblock, eblock, ecblock and fblock instead.
func GetBlockByHeightRaw(blockType string, height int64) (*BlockByHeightRawResponse, error) {
	params := heightRequest{Height: height, NoRaw: false} // include raw
	req := NewJSON2Request(fmt.Sprintf("%vblock-by-height", blockType), APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	block := new(BlockByHeightRawResponse)
	if err := json.Unmarshal(resp.JSONResult(), block); err != nil {
		return nil, err
	}

	return block, nil
}
