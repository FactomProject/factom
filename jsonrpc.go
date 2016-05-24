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

func EncodeJSON(data interface{}) ([]byte, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := EncodeJSON(data)
	if err != nil {
		return "", err
	}
	return string(encoded), err
}

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewJSONError(code int, message string, data interface{}) *JSONError {
	j := new(JSONError)
	j.Code = code
	j.Message = message
	j.Data = data
	return j
}

func (e *JSONError) Error() string {
	s := fmt.Sprint(e.Message)
	if e.Data != nil {
		s += fmt.Sprint(": ", e.Data)
	}
	return s
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

func ParseJSON2Request(request string) (*JSON2Request, error) {
	j := new(JSON2Request)
	err := json.Unmarshal([]byte(request), j)
	if err != nil {
		return nil, err
	}
	if j.JSONRPC != "2.0" {
		return nil, fmt.Errorf("Invalid JSON RPC version - `%v`, should be `2.0`", j.JSONRPC)
	}
	return j, nil
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

func (j *JSON2Response) JSONString() (string, error) {
	return EncodeJSONString(j)
}

func (j *JSON2Response) JSONResult() []byte {
	p, _ := json.Marshal(j.Result)
	return p
}

func (j *JSON2Response) String() string {
	str, _ := j.JSONString()
	return str
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

func walletRequest(req *JSON2Request) (*JSON2Response, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/v2", walletServer),
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

var apiCounter = newCounter()

//Requests

type NameRequest struct {
	Name string `json:"name"`
}

type AddressRequest struct {
	Address string `json:"address"`
}

type ChainIDRequest struct {
	ChainID string `json:"chainid"`
}

type EntryRequest struct {
	Entry string `json:"entry"`
}

type HashRequest struct {
	Hash string `json:"hash"`
}

type KeyMRRequest struct {
	KeyMR string `json:"keymr"`
}

type KeyRequest struct {
	Key string `json:"key"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

type TransactionRequest struct {
	Transaction string `json:"transaction"`
}
