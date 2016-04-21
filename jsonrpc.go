// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JSON2Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Params  interface{} `json:"params,omitempty"`
	Method  string      `json:"method,omitempty"`
}

func NewJSON2Request(method string, id, params interface{}) *JSON2Request {
	j := new(JSON2Request)
	j.JSONRPC = "2.0"
	j.ID = id
	j.Params = params
	j.Method = method
	return j
}

type JSON2Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Error   *JSONError  `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func NewJSON2Response() *JSON2Response {
	j := new(JSON2Response)
	j.JSONRPC = "2.0"
	return j
}

func factomdRequest(req *JSON2Request) (*JSON2Response, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v2", factomdServer),
		"application/json",
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := NewJSON2Response()
	if err := json.Unmarshal(body, r); err != nil {
		return nil, err
	}

	return r, nil
}

// newCounter is used to generate the ID field for the JSON2Request
func newCounter() func() int {
	count := 0
	return func() int {
		count += 1
		return count
	}
}
