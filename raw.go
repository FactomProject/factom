package factom

import (
	"encoding/hex"
	"fmt"

	"github.com/FactomProject/factomd/wsapi"
)

func GetRaw(keymr string) ([]byte, error) {
	resp, err := CallV2("get-raw-data", false, keymr)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	raw, err := hex.DecodeString(resp.Result.(*wsapi.GetRawDataResponse).Data)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
