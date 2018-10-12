package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
)

// An Identity is an array of names and a hierarchy of keys. It can assign/receive
// Attributes as JSON objects and rotate/replace its currently valid keys.
type Identity struct {
	ChainID string
	Name    []string
	Keys    []*IdentityKey
}

type IdentityAttribute struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

// GetIdentityChainID takes an identity name and returns its corresponding ChainID
func GetIdentityChainID(name []string) string {
	hs := sha256.New()
	for _, part := range name {
		h := sha256.Sum256([]byte(part))
		hs.Write(h[:])
	}
	return hex.EncodeToString(hs.Sum(nil))
}

// NewIdentityChain creates an returns a Chain struct for a new identity. Publish it to the
// blockchain using the usual factom.CommitChain(...) and factom.RevealChain(...) calls.
func NewIdentityChain(name []string, keys []*IdentityKey) *Chain {
	e := &Entry{}
	for _, part := range name {
		e.ExtIDs = append(e.ExtIDs, []byte(part))
	}

	var publicKeys []string
	for _, key := range keys {
		publicKeys = append(publicKeys, key.PubString())
	}
	keysMap := map[string][]string{"keys": publicKeys}
	keysJSON, _ := json.Marshal(keysMap)
	e.Content = keysJSON
	c := NewChain(e)
	return c
}

// GetKeysAtHeight returns the identity's public keys that were/are valid at the highest saved block height
func (i *Identity) GetKeysAtCurrentHeight() ([]*IdentityKey, error) {
	heights, err := GetHeights()
	if err != nil {
		return nil, err
	}
	return i.GetKeysAtHeight(heights.DirectoryBlockHeight)
}

// GetKeysAtHeight returns the identity's public keys that were valid at the specified block height
func (i *Identity) GetKeysAtHeight(height int64) ([]*IdentityKey, error) {
	entries, err := GetAllChainEntriesAtHeight(i.ChainID, height)
	if err != nil {
		return nil, err
	}

	var initialKeys map[string][]string
	initialKeysJSON := entries[0].Content
	err = json.Unmarshal(initialKeysJSON, &initialKeys)
	if err != nil {
		fmt.Println("Failed to unmarshal json from initial key declaration")
		return nil, err
	}

	var validKeys []*IdentityKey
	for _, pubString := range initialKeys["keys"] {
		if IdentityKeyStringType(pubString) != IDPub {
			return nil, fmt.Errorf("invalid Identity Public Key string in first entry")
		}
		pub := base58.Decode(pubString)
		k := NewIdentityKey()
		copy(k.Pub[:], pub[IDKeyPrefixLength:IDKeyBodyLength])
		validKeys = append(validKeys, k)
	}

	for _, e := range entries {
		if len(e.ExtIDs) < 5 || bytes.Compare(e.ExtIDs[0], []byte("ReplaceKey")) != 0 {
			continue
		}
		if len(e.ExtIDs[1]) != 55 || len(e.ExtIDs[2]) != 55 || len(e.ExtIDs[3]) != 64 {
			continue
		}

		var oldKey [32]byte
		oldPubString := string(e.ExtIDs[1])
		if IdentityKeyStringType(oldPubString) != IDPub {
			continue
		}
		b := base58.Decode(oldPubString)
		copy(oldKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

		var newKey [32]byte
		newPubString := string(e.ExtIDs[2])
		if IdentityKeyStringType(newPubString) != IDPub {
			continue
		}
		b = base58.Decode(newPubString)
		copy(newKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

		var signature [64]byte
		copy(signature[:], e.ExtIDs[3])
		signerPubString := string(e.ExtIDs[4])

		levelToReplace := -1
		for level, key := range validKeys {
			if bytes.Compare(oldKey[:], key.PubBytes()) == 0 {
				levelToReplace = level
			}
		}
		if levelToReplace == -1 {
			// oldkey not in the set of valid keys when this entry was published
			continue
		}

		message := []byte(oldPubString + newPubString)
		for level, key := range validKeys {
			if level > levelToReplace {
				// low priority key trying to replace high priority key, disregard
				break
			}
			if key.PubString() == signerPubString && ed.Verify(key.Pub, message, &signature) {
				validKeys[levelToReplace].Pub = &newKey
				break
			}
		}
	}
	return validKeys, nil
}

// NewIdentityKeyReplacementEntry creates and returns a new Entry struct for the key replacement. Publish it to the
// blockchain using the usual factom.CommitEntry(...) and factom.RevealEntry(...) calls.
func NewIdentityKeyReplacementEntry(chainID string, oldKey *IdentityKey, newKey *IdentityKey, signerKey *IdentityKey) *Entry {
	message := []byte(oldKey.String() + newKey.String())
	signature := signerKey.Sign(message)

	e := Entry{}
	e.ChainID = chainID
	e.ExtIDs = [][]byte{[]byte("ReplaceKey"), []byte(oldKey.String()), []byte(newKey.String()), signature[:], []byte(signerKey.String())}
	return &e
}

// NewIdentityAttributeEntry creates and returns an Entry struct that assigns an attribute JSON object to a given
// identity. Publish it to the blockchain using the usual factom.CommitEntry(...) and factom.RevealEntry(...) calls.
func NewIdentityAttributeEntry(receiverChainID string, destinationChainID string, attributesJSON string, signerKey *IdentityKey, signerChainID string) *Entry {
	message := []byte(receiverChainID + destinationChainID)
	attributeHash := sha256.Sum256([]byte(attributesJSON))
	message = append(message, attributeHash[:]...)
	signature := signerKey.Sign(message)

	e := Entry{}
	e.ChainID = destinationChainID
	e.ExtIDs = [][]byte{[]byte("IdentityAttribute"), []byte(receiverChainID), signature[:], []byte(signerKey.String()), []byte(signerChainID)}
	e.Content = []byte(attributesJSON)
	return &e
}

// NewIdentityAttributeEndorsementEntry creates and returns an Entry struct that agrees with or recognizes a given
// attribute. Publish it to the blockchain using the usual factom.CommitEntry(...) and factom.RevealEntry(...) calls.
func NewIdentityAttributeEndorsementEntry(destinationChainID string, attributeEntryHash string, signerKey *IdentityKey, signerChainID string) *Entry {
	message := []byte(destinationChainID + attributeEntryHash)
	signature := signerKey.Sign(message)

	e := Entry{}
	e.ChainID = destinationChainID
	e.ExtIDs = [][]byte{[]byte("IdentityAttributeEndorsement"), signature[:], []byte(signerKey.String()), []byte(signerChainID)}
	e.Content = []byte(attributeEntryHash)
	return &e
}

// IsValidAttribute returns true if the entry is a properly formatted attribute with a verifiable signature.
// Note: does not check that the signer key was valid for the signer identity at the time of publishing.
func IsValidAttribute(e *Entry) bool {
	// Check ExtIDs for valid formatting, then process them
	if len(e.ExtIDs) < 5 || bytes.Compare(e.ExtIDs[0], []byte("IdentityAttribute")) != 0 {
		return false
	}
	receiverChainID := string(e.ExtIDs[1])
	signerChainID := string(e.ExtIDs[4])
	if len(receiverChainID) != 64 || len(signerChainID) != 64 {
		return false
	}
	var signature [64]byte
	copy(signature[:], e.ExtIDs[2])
	var signerKey [32]byte
	signerPubString := string(e.ExtIDs[3])
	if IdentityKeyStringType(signerPubString) != IDPub {
		return false
	}
	b := base58.Decode(signerPubString)
	copy(signerKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

	// Message that was signed = ReceiverChainID + DestinationChainID + AttributesJSON
	msg := []byte(receiverChainID + e.ChainID)
	attributesHash := sha256.Sum256(e.Content)
	msg = append(msg, attributesHash[:]...)
	return ed.Verify(&signerKey, msg, &signature)
}

// IsValidEndorsement returns true if the Entry is a properly formatted attribute endorsement with a verifiable signature.
// Note: does not check that the signer key was valid for the signer identity at the time of publishing.
func IsValidEndorsement(e *Entry) bool {
	// Check ExtIDs for valid formatting, then process them
	if len(e.ExtIDs) < 4 || string(e.ExtIDs[0]) != "IdentityAttributeEndorsement" {
		return false
	}

	signerChainID := string(e.ExtIDs[3])
	if len(signerChainID) != 64 {
		return false
	}
	var signature [64]byte
	copy(signature[:], e.ExtIDs[1])
	var signerKey [32]byte
	signerPubString := string(e.ExtIDs[2])
	if IdentityKeyStringType(signerPubString) != IDPub {
		return false
	}
	b := base58.Decode(signerPubString)
	copy(signerKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

	// Message that was signed = DestinationChainID + AttributeEntryHash
	msg := e.ChainID + string(e.Content)
	return ed.Verify(&signerKey, []byte(msg), &signature)
}
