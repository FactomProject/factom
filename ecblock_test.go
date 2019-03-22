// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"github.com/FactomProject/factom"
)

func TestGetECBlock(t *testing.T) {
	ecb, err := factom.GetECBlock("c3836152f2f28c55f0b31807eb0c1ef8d1a4a16241b8ca4612313c75bf38a541")
	if err != nil {
		t.Error(err)
	}
	t.Log(ecb)
}
