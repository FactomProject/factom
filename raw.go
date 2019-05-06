// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"encoding/hex"
	"encoding/json"
)

type RawData struct {
	Data string `json:"data"`
}

func (r *RawData) GetDataBytes() ([]byte, error) {
	return hex.DecodeString(r.Data)
}

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
