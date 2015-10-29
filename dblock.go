package factom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetDBlockHeight() (int, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://%s/v1/directory-block-height/", server))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf(string(body))
	}
	type dbh struct {
		Height int
	}
	d := new(dbh)
	if err := json.Unmarshal(body, d); err != nil {
		return 0, fmt.Errorf("%s: %s\n", err, body)
	}

	return d.Height, nil
}

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

func (d *DBlock) String() string {
	var s string
	s += fmt.Sprintln("PrevBlockKeyMR:", d.Header.PrevBlockKeyMR)
	s += fmt.Sprintln("Timestamp:", d.Header.Timestamp)
	s += fmt.Sprintln("SequenceNumber:", d.Header.SequenceNumber)
	for _, v := range d.EntryBlockList {
		s += fmt.Sprintln("EntryBlock {")
		s += fmt.Sprintln("	ChainID", v.ChainID)
		s += fmt.Sprintln("	KeyMR", v.KeyMR)
		s += fmt.Sprintln("}")
	}
	return s
}
