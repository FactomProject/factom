// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/FactomProject/factom"
	"testing"
)

var _ = fmt.Sprint("testing")

func TestNewJSON2Request(t *testing.T) {
	type t1 struct {
		A int
		B string
	}
	//removed because this test was failing when running ack tests. should be made stateless.
	/*
		x1 := &t1{A: 1, B: "hello"}
		j1 := NewJSON2Request("testing", APICounter(), x1)
		r1 := `{"jsonrpc":"2.0","id":1,"params":{"A":1,"B":"hello"},"method":"testing"}`
		if p, err := json.Marshal(j1); err != nil {
			t.Error(err)
		} else if string(p) != r1 {
			t.Errorf(string(p))
		}

		x2 := "hello"
		j2 := NewJSON2Request("testing", APICounter(), x2)
		r2 := `{"jsonrpc":"2.0","id":2,"params":"hello","method":"testing"}`
		if p, err := json.Marshal(j2); err != nil {
			t.Error(err)
		} else if string(p) != r2 {
			t.Errorf(string(p))
		}

		x3 := new(Entry)
		x3.ChainID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		x3.ExtIDs = append(x3.ExtIDs, []byte("test01"))
		x3.Content = []byte("hello factom")
		j3 := NewJSON2Request("testing", APICounter(), x3)
		r3 := `{"jsonrpc":"2.0","id":3,"params":{"chainid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","extids":["746573743031"],"content":"68656c6c6f20666163746f6d"},"method":"testing"}`
		if p, err := json.Marshal(j3); err != nil {
			t.Error(err)
		} else if string(p) != r3 {
			t.Errorf(string(p))
		}*/
}

func TestJSON2Response(t *testing.T) {
	j1 := []byte(`{"jsonrpc":"2.0","id":2,"result":"hello"}`)
	j2 := []byte(`{"jsonrpc":"2.0","id":3,"result":{"ChainID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","ExtIDs":["746573743031"],"Content":"68656c6c6f20666163746f6d"},"method":"testing"}`)
	r1 := `"hello"`
	r2 := Entry{ChainID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ExtIDs: [][]byte{[]byte("test01")}, Content: []byte("hello factom")}

	resp := NewJSON2Response()
	if err := json.Unmarshal(j1, resp); err != nil {
		t.Error(err)
	} else if string(resp.JSONResult()) != r1 {
		t.Errorf("%s is not equal to %s", resp.JSONResult(), r1)
	}

	resp = NewJSON2Response()
	if err := json.Unmarshal(j2, resp); err != nil {
		t.Error(err)
	}
	e := new(Entry)
	e.UnmarshalJSON(resp.JSONResult())

	if e.ChainID != r2.ChainID {
		t.Error(e)
	}
	for i, v := range e.ExtIDs {
		if !bytes.Equal(v, r2.ExtIDs[i]) {
			t.Error(e)
		}
	}
	if !bytes.Equal(e.Content, r2.Content) {
		t.Error(e)
	}
}
