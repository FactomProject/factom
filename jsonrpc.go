// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RPCConfig struct {
	WalletTLSEnable    bool
	WalletTLSKeyFile   string
	WalletTLSCertFile  string
	WalletRPCUser      string
	WalletRPCPassword  string
	FactomdTLSEnable   bool
	FactomdTLSCertFile string
	FactomdRPCUser     string
	FactomdRPCPassword string
	FactomdServer      string
	WalletServer       string
	WalletTimeout      time.Duration
	FactomdTimeout     time.Duration
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

func (j *JSON2Request) JSONString() (string, error) {
	return EncodeJSONString(j)
}

func (j *JSON2Request) String() string {
	str, _ := j.JSONString()
	return str
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

func SetFactomdEncryption(tls bool, certFile string) {
	RpcConfig.FactomdTLSEnable = tls
	RpcConfig.FactomdTLSCertFile = certFile
}

func GetFactomdEncryption() (bool, string) {
	return RpcConfig.FactomdTLSEnable, RpcConfig.FactomdTLSCertFile
}

func SetFactomdTimeout(timeout time.Duration) {
	RpcConfig.FactomdTimeout = timeout
}

func GetFactomdTimeout() time.Duration {
	return RpcConfig.FactomdTimeout
}

func SetWalletTimeout(timeout time.Duration) {
	RpcConfig.WalletTimeout = timeout
}

func GetWalletTimeout() time.Duration {
	return RpcConfig.WalletTimeout
}

func SetWalletRpcConfig(user string, password string) {
	RpcConfig.WalletRPCUser = user
	RpcConfig.WalletRPCPassword = password
}

func GetWalletRpcConfig() (string, string) {
	return RpcConfig.WalletRPCUser, RpcConfig.WalletRPCPassword
}

func SetWalletEncryption(tls bool, certFile string) {
	RpcConfig.WalletTLSEnable = tls
	RpcConfig.WalletTLSCertFile = certFile
}

func GetWalletEncryption() (bool, string) {
	return RpcConfig.WalletTLSEnable, RpcConfig.WalletTLSCertFile
}

// SetFactomdServer sets where to find the factomd server, and tells the server its public ip
func SetFactomdServer(s string) {
	RpcConfig.FactomdServer = s
}

// SetWalletServer sets where to find the fctwallet server, and tells the server its public ip
func SetWalletServer(s string) {
	RpcConfig.WalletServer = s
}

// FactomdServer returns where to find the factomd server, and tells the server its public ip
func FactomdServer() string {
	return RpcConfig.FactomdServer
}

// FactomdServer returns where to find the fctwallet server, and tells the server its public ip
func WalletServer() string {
	return RpcConfig.WalletServer
}

// SendFactomdRequest sends a json object to factomd
func SendFactomdRequest(req *JSON2Request) (*JSON2Response, error) {
	return factomdRequest(req)
}

func factomdRequest(req *JSON2Request) (*JSON2Response, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	factomdTls, factomdCertPath := GetFactomdEncryption()

	var client *http.Client
	var httpx string

	if factomdTls == true {
		caCert, err := ioutil.ReadFile(factomdCertPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tr := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}}

		client = &http.Client{Transport: tr, Timeout: GetFactomdTimeout()}
		httpx = "https"

	} else {
		client = &http.Client{Timeout: GetFactomdTimeout()}
		httpx = "http"
	}
	re, err := http.NewRequest("POST",
		fmt.Sprintf("%s://%s/v2", httpx, RpcConfig.FactomdServer),
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	user, pass := GetFactomdRpcConfig()
	re.SetBasicAuth(user, pass)
	re.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(re)
	if err != nil {
		errs := fmt.Sprintf("%s", err)
		if strings.Contains(errs, "\\x15\\x03\\x01\\x00\\x02\\x02\\x16") {
			err = fmt.Errorf("Factomd API connection is encrypted. Please specify -factomdtls=true and -factomdcert=factomdAPIpub.cert (%v)", err.Error())
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Factomd username/password incorrect.  Edit factomd.conf or\ncall factom-cli with -factomduser=<user> -factomdpassword=<pass>")
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

	walletTls, walletCertPath := GetWalletEncryption()

	var client *http.Client
	var httpx string

	if walletTls == true {
		caCert, err := ioutil.ReadFile(walletCertPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tr := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}}

		client = &http.Client{Transport: tr, Timeout: GetWalletTimeout()}
		httpx = "https"

	} else {
		client = &http.Client{Timeout: GetWalletTimeout()}
		httpx = "http"
	}

	re, err := http.NewRequest("POST",
		fmt.Sprintf("%s://%s/v2", httpx, RpcConfig.WalletServer),
		bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	user, pass := GetWalletRpcConfig()
	re.SetBasicAuth(user, pass)
	re.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(re)
	if err != nil {
		errs := fmt.Sprintf("%s", err)
		if strings.Contains(errs, "\\x15\\x03\\x01\\x00\\x02\\x02\\x16") {
			err = fmt.Errorf("Factom-walletd API connection is encrypted. Please specify -wallettls=true and -walletcert=walletAPIpub.cert (%v)", err.Error())
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Wallet username/password incorrect.  Edit factomd.conf or\ncall factom-cli with -walletuser=<user> -walletpassword=<pass>")
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
