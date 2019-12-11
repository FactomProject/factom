// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// CurrentMinuteInfo represents the current state of the factom network from the
// factomd API.
type CurrentMinuteInfo struct {
	LeaderHeight            int64 `json:"leaderheight"`
	DirectoryBlockHeight    int64 `json:"directoryblockheight"`
	Minute                  int64 `json:"minute"`
	CurrentBlockStartTime   int64 `json:"currentblockstarttime"`
	CurrentMinuteStartTime  int64 `json:"currentminutestarttime"`
	CurrentTime             int64 `json:"currenttime"`
	DirectoryBlockInSeconds int64 `json:"directoryblockinseconds"`
	StallDetected           bool  `json:"stalldetected"`
	FaultTimeout            int64 `json:"faulttimeout"`
	RoundTimeout            int64 `json:"roundtimeout"`
}

func (c *CurrentMinuteInfo) String() string {
	var s string

	s += fmt.Sprintln("LeaderHeight:", c.LeaderHeight)
	s += fmt.Sprintln("DirectoryBlockHeight:", c.DirectoryBlockHeight)
	s += fmt.Sprintln("Minute:", c.Minute)
	s += fmt.Sprintln("CurrentBlockStartTime:", c.CurrentBlockStartTime)
	s += fmt.Sprintln("CurrentMinuteStartTime:", c.CurrentMinuteStartTime)
	s += fmt.Sprintln("CurrentTime:", c.CurrentTime)
	s += fmt.Sprintln("DirectoryBlockInSeconds:", c.DirectoryBlockInSeconds)
	s += fmt.Sprintln("StallDetected:", c.StallDetected)
	s += fmt.Sprintln("FaultTimeout:", c.FaultTimeout)
	s += fmt.Sprintln("RoundTimeout:", c.RoundTimeout)

	return s
}

// GetCurrentMinute gets the current network information from the factom daemon.
func GetCurrentMinute() (*CurrentMinuteInfo, error) {
	req := NewJSON2Request("current-minute", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	c := new(CurrentMinuteInfo)
	if err := json.Unmarshal(resp.JSONResult(), c); err != nil {
		return nil, err
	}

	return c, nil
}
