// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	. "github.com/FactomProject/factom"
)

func TestGetMultipleFCTBalances(t *testing.T) {
	badfa := "abcdef"
	if bs, err := GetMultipleFCTBalances(badfa); err != nil {
		t.Error(err)
	} else if bs.Balances[0].Err != "Error decoding address" {
		t.Error("should have recieved error for bad address instead got", err)
	}
	fas := []string{
		"FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu",
		"FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu",
		"FA3upjWMKHmStAHR5ZgKVK4zVHPb8U74L2wzKaaSDQEonHajiLeq",
	}
	bs, err := GetMultipleFCTBalances(fas...)
	if err != nil {
		t.Error(err)
	}
	t.Log(bs)
}
