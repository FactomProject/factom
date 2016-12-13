// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	. "github.com/FactomProject/factom"
	"testing"
)

func TestFactoidToFactoshi(t *testing.T) {
	v1 := "1.0001"
	v2 := "100000000.00000001"
	v3 := ".01"
	
	e1 := uint64(100010000)
	e2 := uint64(10000000000000001)
	e3 := uint64(1000000)
	
	if r1 := FactoidToFactoshi(v1); r1 != e1 {
		t.Errorf("r1=%d expecting %d", r1, e1)
	}
	
	if r2 := FactoidToFactoshi(v2); r2 != e2 {
		t.Errorf("r2=%d expecting %d", r2, e2)
	}

	if r3 := FactoidToFactoshi(v3); r3 != e3 {
		t.Errorf("r3=%d expecting %d", r3, e3)
	}
}
