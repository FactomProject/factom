// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"fmt"
)

// Diagnostics represents a set of diagnostic/debugging information about
// factomd and the Factom Network.
type Diagnostics struct {
	Name                  string `json:"name"`
	ID                    string `json:"id,omitempty"`
	PublicKey             string `json:"publickey,omitempty"`
	Role                  string `json:"role"`
	LeaderHeight          int    `json:"leaderheight"`
	CurrentMinute         int    `json:"currentminute"`
	CurrentMinuteDuration int64  `json:"currentminuteduration"`
	PrevMinuteDuration    int64  `json:"previousminuteduration"`
	BalanceHash           string `json:"balancehash"`
	TempBalanceHash       string `json:"tempbalancehash"`
	LastBlockFromDBState  bool   `json:"lastblockfromdbstate"`

	SyncInfo struct {
		Status   string   `json:"status"`
		Received int      `json:"received,omitempty"`
		Expected int      `json:"expected,omitempty"`
		Missing  []string `json:"missing,omitempty"`
	} `json:"syncing"`

	AuthSet struct {
		Leaders []struct {
			ID                string `json:"id"`
			VM                int    `json:"vm"`
			ProcessListHeight int    `json:"listheight"`
			ListLength        int    `json:"listlength"`
			NextNil           int    `json:"nextnil"`
		} `json:"leaders"`

		Audits []struct {
			ID     string `json:"id"`
			Online bool   `json:"online"`
		} `json:"audits"`
	} `json:"authset"`

	ElectionInfo struct {
		InProgress bool   `json:"inprogress"`
		VMIndex    int    `json:"vmindex,omitempty"`
		FedIndex   int    `json:"fedindex,omitempty"`
		FedID      string `json:"fedid,omitempty"`
		Round      int    `json:"round,omitempty"`
	} `json:"elections"`
}

func (d *Diagnostics) String() string {
	var s string

	s += fmt.Sprintln("Name:", d.Name)
	s += fmt.Sprintln("ID:", d.ID)
	s += fmt.Sprintln("PublicKey:", d.PublicKey)
	s += fmt.Sprintln("Role:", d.Role)
	s += fmt.Sprintln("LeaderHeight:", d.LeaderHeight)
	s += fmt.Sprintln("CurrentMinute:", d.CurrentMinute)
	s += fmt.Sprintln("CurrentMinuteDuration:", d.CurrentMinuteDuration)
	s += fmt.Sprintln("PrevMinuteDuration:", d.PrevMinuteDuration)
	s += fmt.Sprintln("BalanceHash:", d.BalanceHash)
	s += fmt.Sprintln("TempBalanceHash:", d.TempBalanceHash)
	s += fmt.Sprintln("LastBlockFromDBState:", d.LastBlockFromDBState)
	// SyncInfo
	s += fmt.Sprintln("Status:", d.SyncInfo.Status)
	s += fmt.Sprintln("Received:", d.SyncInfo.Received)
	s += fmt.Sprintln("Expected:", d.SyncInfo.Expected)
	for _, m := range d.SyncInfo.Missing {
		s += fmt.Sprintln("Missing:", m)
	}
	// ElectionInfo
	s += fmt.Sprintln("InProgress:", d.ElectionInfo.InProgress)
	s += fmt.Sprintln("VMIndex:", d.ElectionInfo.VMIndex)
	s += fmt.Sprintln("FedIndex:", d.ElectionInfo.FedIndex)
	s += fmt.Sprintln("FedID:", d.ElectionInfo.FedID)
	s += fmt.Sprintln("Round:", d.ElectionInfo.Round)
	// AuthSet
	s += fmt.Sprintln("Leaders {")
	for _, v := range d.AuthSet.Leaders {
		s += fmt.Sprintln(" ID:", v.ID)
		s += fmt.Sprintln(" VM:", v.VM)
		s += fmt.Sprintln(" ProcessListHeight:", v.ProcessListHeight)
		s += fmt.Sprintln(" ListLength:", v.ListLength)
		s += fmt.Sprintln(" NextNil:", v.NextNil)
	}
	s += fmt.Sprintln("}") // Leaders
	s += fmt.Sprintln("Audits {")
	for _, v := range d.AuthSet.Audits {
		s += fmt.Sprintln(" ID:", v.ID)
		s += fmt.Sprintln(" Online:", v.Online)
	}
	s += fmt.Sprintln("}") // Audits

	return s
}

// GetDiagnostics reads diagnostic information from factomd.
func GetDiagnostics() (*Diagnostics, error) {
	req := NewJSON2Request("diagnostics", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	d := new(Diagnostics)
	if err := json.Unmarshal(resp.JSONResult(), d); err != nil {
		return nil, err
	}

	return d, nil
}
