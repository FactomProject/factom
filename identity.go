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
func NewIdentityChain(name []string, keys []string) (*Chain, error) {
	e := &Entry{}
	e.ExtIDs = append(e.ExtIDs, []byte("IdentityChain"))
	for _, part := range name {
		e.ExtIDs = append(e.ExtIDs, []byte(part))
	}

	var publicKeys []string
	for _, key := range keys {
		if IdentityKeyStringType(key) != IDPub {
			return nil, fmt.Errorf("provided key %s is not a valid identity public key", key)
		}
		publicKeys = append(publicKeys, key)
	}
	keysMap := map[string]interface{}{"version": 1, "keys": publicKeys}
	keysJSON, _ := json.Marshal(keysMap)
	e.Content = keysJSON
	c := NewChain(e)
	return c, nil
}

// GetActiveIdentityKeys returns the identity's public keys that were/are active at the highest saved block height,
// along with that blockheight
func GetActiveIdentityKeys(chainID string) ([]string, int64, error) {
	heights, err := GetHeights()
	if err != nil {
		return nil, -1, err
	}
	keys, err := GetActiveIdentityKeysAtHeight(chainID, heights.DirectoryBlockHeight)
	return keys, heights.DirectoryBlockHeight, err
}

// GetActiveIdentityKeysAtHeight returns the identity's public keys that were active at the specified block height
func GetActiveIdentityKeysAtHeight(chainID string, height int64) ([]string, error) {
	if !ChainExists(chainID) {
		return nil, fmt.Errorf("chain does not exist")
	}

	entries, err := GetAllChainEntriesAtHeight(chainID, height)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, fmt.Errorf("chain did not yet exist at height %d", height)
	} else if len(entries[0].ExtIDs) == 0 || bytes.Compare(entries[0].ExtIDs[0], []byte("IdentityChain")) != 0 {
		return nil, fmt.Errorf("no identity found at chain ID: %s", chainID)
	}

	var identityInfo struct {
		Version     int      `json:"version"`
		InitialKeys []string `json:"keys"`
	}
	initialKeysJSON := entries[0].Content
	err = json.Unmarshal(initialKeysJSON, &identityInfo)
	if err != nil {
		return nil, fmt.Errorf("no identity found at chain ID: %s", chainID)
	}

	var activeKeys []*IdentityKey
	allKeys := make(map[string]bool)
	for _, pubString := range identityInfo.InitialKeys {
		if IdentityKeyStringType(pubString) != IDPub {
			return nil, fmt.Errorf("invalid identity public key string in first entry: %s", pubString)
		} else if _, present := allKeys[pubString]; present {
			continue
		}
		pub := base58.Decode(pubString)
		k := NewIdentityKey()
		copy(k.Pub[:], pub[IDKeyPrefixLength:IDKeyBodyLength])
		activeKeys = append(activeKeys, k)
		allKeys[pubString] = true
	}

	for _, e := range entries {
		if len(e.ExtIDs) < 5 || bytes.Compare(e.ExtIDs[0], []byte("ReplaceKey")) != 0 {
			continue
		}
		if len(e.ExtIDs[1]) != 55 || len(e.ExtIDs[2]) != 55 || len(e.ExtIDs[3]) != ed.SignatureSize {
			continue
		}

		var oldKey [ed.PublicKeySize]byte
		oldPubString := string(e.ExtIDs[1])
		if IdentityKeyStringType(oldPubString) != IDPub {
			continue
		}
		b := base58.Decode(oldPubString)
		copy(oldKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

		var newKey [ed.PublicKeySize]byte
		newPubString := string(e.ExtIDs[2])
		if IdentityKeyStringType(newPubString) != IDPub {
			continue
		}
		b = base58.Decode(newPubString)
		copy(newKey[:], b[IDKeyPrefixLength:IDKeyBodyLength])

		// Disallow re-adding retired or currently active keys
		if _, present := allKeys[newPubString]; present {
			continue
		}

		var signature [ed.SignatureSize]byte
		copy(signature[:], e.ExtIDs[3])
		signerPubString := string(e.ExtIDs[4])

		levelToReplace := -1
		for level, key := range activeKeys {
			if bytes.Compare(oldKey[:], key.PubBytes()) == 0 {
				levelToReplace = level
			}
		}
		if levelToReplace == -1 {
			// oldkey not in the set of valid keys when this entry was published
			continue
		}

		message := []byte(chainID + oldPubString + newPubString)
		for level, key := range activeKeys {
			if level > levelToReplace {
				// low priority key trying to replace high priority key, disregard
				break
			}
			if key.PubString() == signerPubString && ed.Verify(key.Pub, message, &signature) {
				activeKeys[levelToReplace].Pub = &newKey
				allKeys[newPubString] = true
				break
			}
		}
	}

	var resp []string
	for _, k := range activeKeys {
		resp = append(resp, k.PubString())
	}
	return resp, nil
}

// NewIdentityKeyReplacementEntry creates and returns a new Entry struct for the key replacement. Publish it to the
// blockchain using the usual factom.CommitEntry(...) and factom.RevealEntry(...) calls.
func NewIdentityKeyReplacementEntry(chainID string, oldKey string, newKey string, signerKey *IdentityKey) (*Entry, error) {
	if IdentityKeyStringType(oldKey) != IDPub {
		return nil, fmt.Errorf("provided key %s is not a valid identity public key", oldKey)
	}
	if IdentityKeyStringType(newKey) != IDPub {
		return nil, fmt.Errorf("provided key %s is not a valid identity public key", newKey)
	}
	message := []byte(chainID + oldKey + newKey)
	signature := signerKey.Sign(message)

	e := Entry{}
	e.ChainID = chainID
	e.ExtIDs = [][]byte{
		[]byte("ReplaceKey"),
		[]byte(oldKey),
		[]byte(newKey),
		signature[:],
		[]byte(signerKey.String()),
	}
	return &e, nil
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
	e.ExtIDs = [][]byte{
		[]byte("IdentityAttribute"),
		[]byte(receiverChainID),
		signature[:],
		[]byte(signerKey.String()),
		[]byte(signerChainID),
	}
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
	e.ExtIDs = [][]byte{
		[]byte("IdentityAttributeEndorsement"),
		signature[:],
		[]byte(signerKey.String()),
		[]byte(signerChainID),
	}
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
	var signature [ed.SignatureSize]byte
	copy(signature[:], e.ExtIDs[2])
	var signerKey [ed.PublicKeySize]byte
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
	var signature [ed.SignatureSize]byte
	copy(signature[:], e.ExtIDs[1])
	var signerKey [ed.PublicKeySize]byte
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
