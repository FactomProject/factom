// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"fmt"

	. "github.com/FactomProject/factom"
)

// Tests reqire a local factomd node to be running and servier the API!

func TestGetDBlock(t *testing.T) {
	d, raw, err := GetDBlock("cde346e7ed87957edfd68c432c984f35596f29c7d23de6f279351cddecd5dc66")
	if err != nil {
		t.Error(err)
	}
	t.Log("dblock:", d)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
}

func TestGetDBlockByHeight(t *testing.T) {
	d, raw, err := GetDBlockByHeight(100)
	if err != nil {
		t.Error(err)
	}
	t.Log("dblock:", d)
	t.Log(fmt.Sprintf("raw: %x\n", raw))
}
