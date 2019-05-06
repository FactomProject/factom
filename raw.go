// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
	"encoding/json"
)

// RawData is a simple hex encoded byte string
type RawData struct {
	Data string `json:"data"`
}

func (r *RawData) GetDataBytes() ([]byte, error) {
	return hex.DecodeString(r.Data)
}

// GetRaw requests the raw data for any binary block kept in the factomd
// database.
func GetRaw(keymr string) ([]byte, error) {
	params := hashRequest{Hash: keymr}
	req := NewJSON2Request("raw-data", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	raw := new(RawData)
	if err := json.Unmarshal(resp.JSONResult(), raw); err != nil {
		return nil, err
	}

	return raw.GetDataBytes()
}

// SendRawMsg sends a raw hex encoded byte string for factomd to send as a
// binary message on the Factom Netwrork.
func SendRawMsg(message string) (string, error) {
	param := messageRequest{Message: message}
	req := NewJSON2Request("send-raw-message", APICounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}

	status := new(struct {
		Message string `json:"message"`
	})
	if err := json.Unmarshal(resp.JSONResult(), status); err != nil {
		return "", err
	}

	return status.Message, nil
}
