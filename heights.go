// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

type HeightsResponse struct {
	DirectoryBlockHeight int64 `json:"directoryblockheight"`
	LeaderHeight         int64 `json:"leaderheight"`
	EntryBlockHeight     int64 `json:"entryblockheight"`
	EntryHeight          int64 `json:"entryheight"`
}

func (d *HeightsResponse) String() string {
	var s string

	s += fmt.Sprintln("DirectoryBlockHeight:", d.DirectoryBlockHeight)
	s += fmt.Sprintln("LeaderHeight:", d.LeaderHeight)
	s += fmt.Sprintln("EntryBlockHeight:", d.EntryBlockHeight)
	s += fmt.Sprintln("EntryHeight:", d.EntryHeight)

	return s
}

func GetHeights() (*HeightsResponse, error) {
	req := NewJSON2Request("heights", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	heights := new(HeightsResponse)
	if err := json.Unmarshal(resp.JSONResult(), heights); err != nil {
		return nil, err
	}

	return heights, nil
}
