package factom

import (
	"encoding/json"
)

type SendRawMessageResponse struct {
	Message string `json:"message"`
}

func SendRawMsg(message string) (*SendRawMessageResponse, error) {
	param := messageRequest{Message: message}
	req := NewJSON2Request("send-raw-message", APICounter(), param)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	status := new(SendRawMessageResponse)
	if err := json.Unmarshal(resp.JSONResult(), status); err != nil {
		return nil, err
	}

	return status, nil
}
