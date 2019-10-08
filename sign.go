package factom

import (
	"encoding/json"
)

type Signature struct {
	PubKey    []byte `json:"pubkey"`
	Signature []byte `json:"signature"`
}

func SignData(addr string, data []byte) (*Signature, error) {
	params := &struct {
		Address string `json:"address"`
		Data    []byte `json:"data"`
	}{
		Address: addr,
		Data:    data,
	}

	req := NewJSON2Request("sign-data", APICounter(), params)
	resp, err := walletRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	sig := new(Signature)
	if err := json.Unmarshal(resp.Result, sig); err != nil {
		return nil, err
	}
	return sig, nil

}
