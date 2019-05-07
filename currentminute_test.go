// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	. "github.com/FactomProject/factom"
)

// TestGetCurrentMinute relies on having a running factom daemon to provide an
// api endpoint at localhost:8088
func TestGetCurrentMinute(t *testing.T) {
	min, err := GetCurrentMinute()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(min.String())
}
