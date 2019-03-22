// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"testing"

	"github.com/FactomProject/factom"
)

func TestGetDiagnostics(t *testing.T) {
	d, err := factom.GetDiagnostics()
	if err != nil {
		t.Error(err)
	}
	t.Log(d)
}
