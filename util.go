// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/FactomProject/factomd/common/primitives"
)

const (
	ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"
)

var (
	server    = "localhost:8088"
	serverFct = "localhost:8089"
)

// SetServer sets the gloabal target for the factomd server
func SetServer(s string) {
	server = s
}

// SetWallet sets the global target for the fctwallet server
func SetWallet(s string) {
	serverFct = s
}

// Server() returns the global server string for debugging
func Server() string {
	return server
}

// milliTime returns a 6 byte slice representing the unix time in milliseconds
func milliTime() (r []byte) {
	buf := new(bytes.Buffer)
	t := time.Now().UnixNano()
	m := t / 1e6
	binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes()[2:]
}

// shad Double Sha256 Hash; sha256(sha256(data))
func shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// sha52
func sha52(data []byte) []byte {
	h1 := sha512.Sum512(data)
	h2 := sha256.Sum256(append(h1[:], data...))
	return h2[:]
}

func CallV2(method string, post bool, params interface{}, dst interface{}) (*primitives.JSON2Response, error) {
	j := primitives.NewJSON2RequestBlank()
	j.Method = method
	j.Params = params
	j.ID = 1

	postGet := "GET"
	if post == true {
		postGet = "POST"
	}

	address := fmt.Sprintf("http://%s/v2", server)

	data, err := j.JSONString()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(postGet, address, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jResp := new(primitives.JSON2Response)
	jResp.Result = dst

	//fmt.Printf("resp body - %v\n", string(body))

	err = json.Unmarshal(body, jResp)
	if err != nil {
		return nil, err
	}

	return jResp, nil
}
