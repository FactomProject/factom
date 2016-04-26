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
	j1 := NewJSON2Request("testing", apiCounter(), x1)
	if p, err := json.Marshal(j1); err != nil {
		t.Error(err)
	} else {
		t.Log(string(p))
	}
	
	x2 := "hello"
	j2 := NewJSON2Request("testing", apiCounter(), x2)
	if p, err := json.Marshal(j2); err != nil {
		t.Error(err)
	} else {
		t.Log(string(p))
	}
	
	x3 := new(Entry)
	x3.ChainID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	x3.ExtIDs = append(x3.ExtIDs, []byte("test01"))
	x3.Content = []byte("hello factom")
	j3 := NewJSON2Request("testing", apiCounter(), x3)
	if p, err := json.Marshal(j3); err != nil {
		t.Error(err)
	} else {
		t.Log(string(p))
	}
}

func TestJSON2Response(t *testing.T) {
	j1 := []byte(`{"jsonrpc":"2.0","id":1,"result":{"A":1,"B":"hello"}}`)
	j2 := []byte(`{"jsonrpc":"2.0","id":2,"result":"hello"}`)
	j3 := []byte(`{"jsonrpc":"2.0","id":3,"result":{"ChainID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","ExtIDs":["746573743031"],"Content":"68656c6c6f20666163746f6d"},"method":"testing"}`)
	j4 := []byte(`{"jsonrpc":"2.0","id":4,"error": {"code": -4, "message": "Method not found"},"result":{"A":1,"B":"hello"}}`)

	resp := NewJSON2Response()
	if err := json.Unmarshal(j1, resp); err != nil {
		t.Error(err)
	}
	t.Log(resp.Result)

	resp = NewJSON2Response()
	if err := json.Unmarshal(j2, resp); err != nil {
		t.Error(err)
	}
	t.Log(resp.Result)

	resp = NewJSON2Response()
	if err := json.Unmarshal(j3, resp); err != nil {
		t.Error(err)
	}
	e := new(Entry)
	if p, err := json.Marshal(resp.Result); err != nil {
		t.Error(err)
	} else {
		if err := e.UnmarshalJSON(p); err != nil {
			t.Error(err)
		}
	}
	t.Log(e)

	resp = NewJSON2Response()
	if err := json.Unmarshal(j4, resp); err != nil {
		t.Error(err)
	}
	if resp.Error != nil {
		t.Log(resp.Error)
	} else {
		t.Error("Expecting Error but no error found!", resp)
	}
}