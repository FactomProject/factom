package factom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	Data string
}

func GetRaw(keymr string) ([]byte, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/get-raw-data/%s", server, keymr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	d := new(Data)
	if err := json.Unmarshal(body, d); err != nil {
		return nil, err
	}

	raw, err := hex.DecodeString(d.Data)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
