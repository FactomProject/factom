package factom

import (
	"encoding/json"
)

type Signature struct {
	PubKey    []byte `json:"pubkey"`
	Signature []byte `json:"signature"`
}

// SignData lets you sign arbitrary data by the specified signer.
// The signer can be either an FA address, EC address, or Identity.
// Be aware that the data is transmitted to the wallet.
func SignData(signer string, data []byte) (*Signature, error) {
	params := &struct {
		Signer string `json:"signer"`
		Data   []byte `json:"data"`
	}{
		Signer: signer,
		Data:   data,
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
