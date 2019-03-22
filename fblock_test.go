// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"github.com/FactomProject/factom"
)

func TestGetTPS(t *testing.T) {
	fb, err := factom.GetFBlock("cfcac07b29ccfa413aeda646b5d386006468189939dfdfa6415b97cc35f2ea1a")
	if err != nil {
		t.Error(err)
	}
	t.Log(fb)
}
