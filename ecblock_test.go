// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"github.com/FactomProject/factom"
)

func TestGetECBlock(t *testing.T) {
	// Check for a missing blockHash
	_, err := factom.GetECBlock("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if err == nil {
		t.Error("expected error for missing block")
	} else {
		t.Log("Missing Block Error:", err)
	}

	ecb, err := factom.GetECBlock("639995e66788ca01709a97684062b466fdce7b840b12861adbe39392f50f6bd3")
	if err != nil {
		t.Error(err)
	}
	t.Log(ecb)
}
