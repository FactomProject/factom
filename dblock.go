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
		TimeStamp      uint64
		SequenceNumber int
	}
	DBEntries []struct {
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	
	d := new(DBlockHead)
	json.Unmarshal(body, d)
	
	return d, nil
}
