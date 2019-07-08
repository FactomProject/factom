// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
)

// GetTPS returns the instant rate (over the previous 3 seconds) and total rate
// (over the lifetime of the node) of Transactions Per Second rate know to
// factomd.
func GetTPS() (instant, total float64, err error) {
	req := NewJSON2Request("tps-rate", APICounter(), nil)
	resp, err := factomdRequest(req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return
	}

	// create temporary type to decode the json tps rate response
	rates := new(struct {
		InstantRate float64 `json:"instanttxrate"`
		TotalRate   float64 `json:"totaltxrate"`
	})

	if err = json.Unmarshal(resp.JSONResult(), rates); err != nil {
		return
	}

	return rates.InstantRate, rates.TotalRate, nil
}
