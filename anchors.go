package factom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Anchors is an anchors response from factomd.
// Note that Ethereum or Bitcoin can be nil
type Anchors struct {
	Height   uint32          `json:"directoryblockheight"`
	KeyMR    string          `json:"directoryblockkeymr"`
	Bitcoin  *AnchorBitcoin  `json:"bitcoin"`
	Ethereum *AnchorEthereum `json:"ethereum"`
}

// AnchorBitcoin is the bitcoin specific anchor
type AnchorBitcoin struct {
	TransactionHash string `json:"transactionhash"`
	BlockHash       string `json:"blockhash"`
}

// AnchorEthereum is the ethereum specific anchor
type AnchorEthereum struct {
	RecordHeight int64        `json:"recordheight"`
	DBHeightMax  int64        `json:"dbheightmax"`
	DBHeightMin  int64        `json:"dbheightmin"`
	WindowMR     string       `json:"windowmr"`
	MerkleBranch []MerkleNode `json:"merklebranch"`

	ContractAddress string `json:"contractaddress"`
	TxID            string `json:"txid"`
	BlockHash       string `json:"blockhash"`
	TxIndex         int64  `json:"txindex"`
}

// MerkleNode is part of the ethereum anchor
type MerkleNode struct {
	Left  string `json:"left,omitempty"`
	Right string `json:"right,omitempty"`
	Top   string `json:"top,omitempty"`
}

func (a *Anchors) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Height: %d\n", a.Height)
	fmt.Fprintf(&sb, "KeyMR: %s\n", a.KeyMR)

	if a.Bitcoin != nil {
		fmt.Fprintf(&sb, "Bitcoin {\n")
		fmt.Fprintf(&sb, " TransactionHash: %s\n", a.Bitcoin.TransactionHash)
		fmt.Fprintf(&sb, " BlockHash: %s\n", a.Bitcoin.BlockHash)
		fmt.Fprintf(&sb, "}\n")
	} else {
		fmt.Fprintf(&sb, "Bitcoin {}\n")
	}

	if a.Ethereum != nil {
		fmt.Fprintf(&sb, "Ethereum {\n")
		fmt.Fprintf(&sb, " RecordHeight: %d\n", a.Ethereum.RecordHeight)
		fmt.Fprintf(&sb, " DBHeightMax: %d\n", a.Ethereum.DBHeightMax)
		fmt.Fprintf(&sb, " DBHeightMin: %d\n", a.Ethereum.DBHeightMin)
		fmt.Fprintf(&sb, " WindowMR: %s\n", a.Ethereum.WindowMR)
		fmt.Fprintf(&sb, "  MerkleBranch {\n")
		for _, branch := range a.Ethereum.MerkleBranch {
			fmt.Fprintf(&sb, "   Branch {\n")
			fmt.Fprintf(&sb, "    Left: %s\n", branch.Left)
			fmt.Fprintf(&sb, "    Right: %s\n", branch.Right)
			fmt.Fprintf(&sb, "    Top: %s\n", branch.Top)
			fmt.Fprintf(&sb, "   Branch }\n")
		}
		fmt.Fprintf(&sb, "  }\n")
		fmt.Fprintf(&sb, " ContractAddress: %s\n", a.Ethereum.ContractAddress)
		fmt.Fprintf(&sb, " TxID: %s\n", a.Ethereum.TxID)
		fmt.Fprintf(&sb, " BlockHash: %s\n", a.Ethereum.BlockHash)
		fmt.Fprintf(&sb, " TxIndex: %d\n", a.Ethereum.TxIndex)
		fmt.Fprintf(&sb, "}\n")
	} else {
		fmt.Fprintf(&sb, "Ethereum {}\n")
	}

	return sb.String()
}

// UnmarshalJSON is an unmarshaller that handles the variable response from factomd
func (a *Anchors) UnmarshalJSON(data []byte) error {
	type tmp *Anchors // unmarshal into a new type to prevent infinite loop
	// json can't unmarshal a bool into a struct, but it can recognize a null pointer
	data = bytes.Replace(data, []byte("\"ethereum\":false"), []byte("\"ethereum\":null"), -1)
	data = bytes.Replace(data, []byte("\"bitcoin\":false"), []byte("\"bitcoin\":null"), -1)
	return json.Unmarshal(data, tmp(a))
}

func getAnchors(hash string, height int64) (*Anchors, error) {
	var params interface{}
	if hash != "" {
		params = hashRequest{Hash: hash}
	} else {
		params = heightRequest{Height: height}
	}
	req := NewJSON2Request("anchors", APICounter(), params)
	resp, err := factomdRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	var res Anchors
	err = json.Unmarshal(resp.Result, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAnchors retrieves the bitcoin and ethereum anchors from factod.
// Hash can be entry hash, entry block keymr, factoid block keymr,
// admin block lookup hash, entry credit block header hash, or
// directory block keymr
func GetAnchors(hash string) (*Anchors, error) {
	return getAnchors(hash, 0)
}

// GetAnchorsByHeight retrieves the bitcoin and ethereum anchors for
// a specific height
func GetAnchorsByHeight(height int64) (*Anchors, error) {
	return getAnchors("", height)
}
