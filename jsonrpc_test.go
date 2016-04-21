// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/json"
	"testing"
)

func TestNewJSON2Request(t *testing.T) {
	type t1 struct {
		A int
		B string
	}
	
	x1 := &t1{A: 1, B: "hello"}
	j1 := NewJSON2Request("testing", counter(), x1)
	if p, err := json.Marshal(j1); err != nil {
		t.Error(err)
	} else {
		t.Log(string(p))
	}
	
	x2 := "hello"
	j2 := NewJSON2Request("testing", counter(), x2)
	if p, err := json.Marshal(j2); err != nil {
		t.Error(err)
	} else {
		t.Log(string(p))
	}
}
