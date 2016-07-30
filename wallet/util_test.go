// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet_test

import (
	. "github.com/FactomProject/factom/wallet"
	"testing"
)

func TestSeedString(t *testing.T) {
	zSeedStr := "sdLGjhUDxGpiBEPRhTwysRYmxNQD6V48Aa84oVzfHvy6suim6qB6m3MCp8aHu1k1CNVLJdB8N9HtGR4NZTtFfp3mj591eA3"

	seed := make([]byte, 64)
	if SeedString(seed) != zSeedStr {
		t.Errorf("seed string does not match")
	}
}
