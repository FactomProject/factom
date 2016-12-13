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
	
	t.Log(FactoidToFactoshi(v1))
	t.Log(FactoidToFactoshi(v2))
	t.Log(FactoidToFactoshi(v3))
}
