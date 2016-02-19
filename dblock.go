package factom

import (
	"fmt"

	"github.com/FactomProject/factomd/wsapi"
)

func GetDBlockHeight() (int, error) {
	resp, err := CallV2("directory-block-height", false, nil, new(wsapi.DirectoryBlockHeightResponse))
	if err != nil {
		return 0, err
	}

	if resp.Error != nil {
		return 0, fmt.Errorf(resp.Error.Message)
	}

	return int(resp.Result.(*wsapi.DirectoryBlockHeightResponse).Height), nil
}

/*
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

*/

func GetDBlock(keymr string) (*wsapi.DirectoryBlockResponse, error) {
	resp, err := CallV2("directory-block-by-keymr", false, keymr, new(wsapi.DirectoryBlockResponse))
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	return resp.Result.(*wsapi.DirectoryBlockResponse), nil
}

func GetDBlockHead() (string, error) {
	resp, err := CallV2("directory-block-head", false, nil, new(wsapi.DirectoryBlockHeadResponse))
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf(resp.Error.Message)
	}

	return resp.Result.(*wsapi.DirectoryBlockHeadResponse).KeyMR, nil
}
