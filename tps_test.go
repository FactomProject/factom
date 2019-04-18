// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	. "github.com/FactomProject/factom"
)

func TestGetTPS(t *testing.T) {
	instant, total, err := GetTPS()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Instant: %f, Total: %f\n", instant, total)
}
