package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DBlock struct {
	DBHash string
	Header struct {
		PrevBlockKeyMR string
		Timestamp      uint64
		SequenceNumber int
	}
	EntryBlockList []struct {
		ChainID string
		KeyMR   string
	}
}

type DBlockHead struct {
	KeyMR string
}

func GetDBlock(keymr string) (*DBlock, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/directory-block-by-keymr/%s", server, keymr))
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

	d := new(DBlock)
	if err := json.Unmarshal(body, d); err != nil {
		return nil, fmt.Errorf("%s: %s\n", err, body)
	}

	return d, nil
}

func GetDBlockHead() (*DBlockHead, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/directory-block-head/", server))
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

	d := new(DBlockHead)
	json.Unmarshal(body, d)

	return d, nil
}
