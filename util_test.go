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
	v4 := "599.9999"
	v5 := "1.99"

	e1 := uint64(100010000)
	e2 := uint64(10000000000000001)
	e3 := uint64(1000000)
	e4 := uint64(59999990000)
	e5 := uint64(199000000)

	if r1 := FactoidToFactoshi(v1); r1 != e1 {
		t.Errorf("r1=%d expecting %d", r1, e1)
	}

	if r2 := FactoidToFactoshi(v2); r2 != e2 {
		t.Errorf("r2=%d expecting %d", r2, e2)
	}

	if r3 := FactoidToFactoshi(v3); r3 != e3 {
		t.Errorf("r3=%d expecting %d", r3, e3)
	}

	if r4 := FactoidToFactoshi(v4); r4 != e4 {
		t.Errorf("r4=%d expecting %d", r4, e4)
	}

	if r5 := FactoidToFactoshi(v5); r5 != e5 {
		t.Errorf("r5=%d expecting %d", r5, e5)
	}
}
