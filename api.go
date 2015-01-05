package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"
	
	"github.com/FactomProject/FactomCode/wallet"
)

var server string = "http://localhost:8083"

func sha(b []byte) []byte {
	s := sha256.New()
	s.Write(b)
	return s.Sum(nil)
}

// PrintEntry is a helper function for debugging entry transport and encoding
func PrintEntry(e *Entry) {
	fmt.Println("ChainID:", hex.EncodeToString(e.ChainID))
	fmt.Println("ExtIDs:")
	for i := range e.ExtIDs {
		fmt.Println("	", string(e.ExtIDs[i]))
	}
	fmt.Println("Data:", string(e.Data))
}

// SetServer specifies the address of the server recieving the factom messages.
// It should be depricated by the final release once the p2p network has been
// implimented
func SetServer(s string) {
	server = s
}

// NewEntry creates a factom entry. It is supplied a string chain id, a []byte
// of data, and a series of string external ids for entry lookup
func NewEntry(cid string, eids []string, data []byte) (e *Entry, err error) {
	e = new(Entry)
	e.ChainID, err = hex.DecodeString(cid)
	if err != nil {
		return nil, err
	}
	e.Data = data
	for _, v := range eids {
		e.ExtIDs = append(e.ExtIDs, []byte(v))
	}
	return
}

// NewChain creates a factom chain from a []string chain name and a new entry
// to be the first entry of the new chain from []byte data, and a series of
// string external ids
func NewChain(name []string, eids []string, data []byte) (c *Chain, err error) {
	c = new(Chain)
	for _, v := range name {
		c.Name = append(c.Name, []byte(v))
	}
	c.GenerateID()
	e := c.FirstEntry
	e.ChainID = c.ChainID
	e.Data = data
	for _, v := range eids {
		e.ExtIDs = append(e.ExtIDs, []byte(v))
	}
	return
}

// CommitEntry sends a message to the factom network containing a hash of the
// entry to be used to verify the later RevealEntry.
func CommitEntry(e *Entry) error {
	var msg bytes.Buffer
	
	msg.Write([]byte{byte(time.Now().Unix())})
	msg.Write(e.MarshalBinary())
	sig := wallet.SignData(msg.Bytes())
	
	// msg.Bytes should be a int64 timestamp followed by a binary entry
	
	data := url.Values{
		"datatype":  {"commitentry"},
		"format":    {"binary"},
		"signature": {hex.EncodeToString((*sig.Sig)[:])},
		"data":      {hex.EncodeToString(msg.Bytes())},
	}
	_, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	return nil
}

// RevealEntry sends a message to the factom network containing the binary
// encoded entry for the server to add it to the factom blockchain. The entry
// will be rejected if a CommitEntry was not done.
func RevealEntry(e *Entry) error {
	data := url.Values{
		"datatype": {"entry"},
		"format":   {"binary"},
		"entry":    {e.Hex()},
	}
	_, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	return nil
}

// CommitChain sends a message to the factom network containing a series of
// hashes to be used to verify the later RevealChain.
func CommitChain(c *Chain) error {
	data := url.Values{
		"datatype": {"chainhash"},
		"format":   {"binary"},
		"data":     {c.Hash()},
	}

	_, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	return nil
}

// RevealChain sends a message to the factom network containing the binary
// encoded first entry for a chain to be used by the server to add a new factom
// chain. It will be rejected if a CommitChain was not done.
func RevealChain(c *Chain) error {
	data := url.Values{
		"datatype": {"entry"},
		"format":   {"binary"},
		"data":     {c.FirstEntry.Hex()},
	}
	_, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	return nil
}

// Submit wraps CommitEntry and RevealEntry. Submit takes a FactomWriter (an
// entry is a FactomWriter) and does a commit and reveal for the entry adding
// it to the factom blockchain.
func Submit(f FactomWriter) (err error) {
	e := f.CreateFactomEntry()
//	err = CommitEntry(e)
//	if err != nil {
//		return err
//	}
	err = RevealEntry(e)
	if err != nil {
		return err
	}
	return nil
}

// CreateChain takes a FactomChainer (a Chain is a FactomChainer) and calls
// commit and reveal to create the factom chain on the network.
func CreateChain(f FactomChainer) error {
	c := f.CreateFactomChain()
	err := CommitChain(c)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Minute)
	err = RevealChain(c)
	if err != nil {
		return err
	}
	return nil
}
