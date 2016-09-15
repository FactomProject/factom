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

type RPCConfig struct {
	TLSEnable          bool   `json:"TLS-enable"`
	TLSKeyFile         string `json:"TLS-keyfile"`
	TLSCertFile        string `json:"TLS-certfile"`
	WalletRPCUser      string `json:"walletrpcuser"`
	WalletRPCPassword  string `json:"walletrpcpassword"`
	FactomdRPCUser     string `json:"factomdrpcuser"`
	FactomdRPCPassword string `json:"factomdrpcpassword"`
	//Authsha     []byte
}

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
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Params  json.RawMessage `json:"params,omitempty"`
	Method  string          `json:"method,omitempty"`
}

func NewJSON2Request(method string, id, params interface{}) *JSON2Request {
	j := new(JSON2Request)
	j.JSONRPC = "2.0"
	j.ID = id
	if b, err := json.Marshal(params); err == nil {
		j.Params = b
	}
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
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Error   *JSONError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
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
	return j.Result
}

func (j *JSON2Response) String() string {
	str, _ := j.JSONString()
	return str
}

func SetFactomdRpcConfig(user string, password string) {
	RpcConfig.FactomdRPCUser = user
	RpcConfig.FactomdRPCPassword = password
}

func GetFactomdRpcConfig() (string, string) {
	return RpcConfig.FactomdRPCUser, RpcConfig.FactomdRPCPassword
}

func SetWalletRpcConfig(user string, password string) {
	RpcConfig.WalletRPCUser = user
	RpcConfig.WalletRPCPassword = password
}

func GetWalletRpcConfig() (string, string) {
	return RpcConfig.WalletRPCUser, RpcConfig.WalletRPCPassword
}

func factomdRequest(req *JSON2Request) (*JSON2Response, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	re, err := http.NewRequest("POST", fmt.Sprintf("http://%s/v2", factomdServer),
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	user, pass := GetFactomdRpcConfig()
	re.SetBasicAuth(user, pass)
	resp, err := client.Do(re)
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

	client := &http.Client{}
	re, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s/v2", walletServer),
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	user, pass := GetWalletRpcConfig()
	re.SetBasicAuth(user, pass)
	re.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(re)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Wallet password protected.  Edit factomd.conf or call factom-cli with -walletuser=<user> -walletpassword=<pass>")
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

var APICounter = newCounter()
