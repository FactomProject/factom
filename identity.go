package factom

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	ed "github.com/FactomProject/ed25519"
)

type Identity struct {
	ChainID string
	Name    []string
	Keys    []*IdentityKey
}

func GetIdentity(name []string) (*Identity, error) {
	// Create the ChainID from a series of hashes of the Identity name
	hs := sha256.New()
	for _, part := range name {
		h := sha256.Sum256([]byte(part))
		hs.Write(h[:])
	}
	chainID := hex.EncodeToString(hs.Sum(nil))
	i := Identity{}
	i.ChainID = chainID
	i.Name = name

	if !ChainExists(chainID) {
		return nil, fmt.Errorf("chain not found")
	}

	keys, err := GetKeysAtCurrentHeight(i.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity keys: %v", err)
	}
	i.Keys = keys
	return &i, nil
}

func CreateIdentityChain(name []string, keys []*IdentityKey, ec *ECAddress) (string, error) {
	e := Entry{}
	for _, part := range name {
		e.ExtIDs = append(e.ExtIDs, []byte(part))
	}

	var publicKeys []string
	for _, key := range keys {
		publicKeys = append(publicKeys, key.PubString())
	}
	keysMap := map[string][]string{"keys": publicKeys}
	keysJSON, _ := json.Marshal(keysMap)
	e.Content = []byte(keysJSON)
	chain := NewChain(&e)

	txID, err := CommitChain(chain, ec)
	if err != nil {
		return "", err
	}
	_, err = RevealChain(chain)
	if err != nil {
		return "", err
	}
	return txID, nil
}

// GetKeysAtHeight returns the identity's public keys that were/are valid at the highest saved block height
func GetKeysAtCurrentHeight(chainID string) ([]*IdentityKey, error) {
	heights, err := GetHeights()
	if err != nil {
		return nil, err
	}
	return GetKeysAtHeight(chainID, heights.DirectoryBlockHeight)
}

// GetKeysAtHeight returns the identity's public keys that were valid at the specified block height
func GetKeysAtHeight(chainID string, height int64) ([]*IdentityKey, error) {
	entries, err := GetAllChainEntriesAtHeight(chainID, height)
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
		pub, err := base64.StdEncoding.DecodeString(pubString)
		if err != nil || len(pub) != 32 {
			return nil, fmt.Errorf("invalid Identity public key string in first entry")
		}
		k := NewIdentityKey()
		copy(k.Pub[:], pub)
		validKeys = append(validKeys, k)
	}

	for _, e := range entries {
		if len(e.ExtIDs) < 4 || bytes.Compare(e.ExtIDs[0], []byte("ReplaceKey")) != 0 {
			continue
		}
		if len(e.ExtIDs[1]) != 32 || len(e.ExtIDs[2]) != 32 || len(e.ExtIDs[3]) != 64 {
			continue
		}
		var oldKey [32]byte
		var newKey [32]byte
		var signature [64]byte
		copy(oldKey[:], e.ExtIDs[1])
		copy(newKey[:], e.ExtIDs[2])
		copy(signature[:], e.ExtIDs[3])

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

		var message []byte
		message = append(message, oldKey[:]...)
		message = append(message, newKey[:]...)
		for level, key := range validKeys {
			if level > levelToReplace {
				// low priority key trying to replace high priority key, disregard
				break
			}
			if ed.Verify(key.Pub, message, &signature) {
				validKeys[levelToReplace].Pub = &newKey
				validKeys[levelToReplace].Sec = new([ed.PrivateKeySize]byte)
				break
			}
		}
	}

	return validKeys, nil
}

func (i *Identity) ReplaceKey(oldKey *[32]byte, newKey *[32]byte, privateKey *[64]byte, ec *ECAddress) (string, error) {
	//publicKey := ed.GetPublicKey(privateKey)
	var message []byte
	message = append(message, oldKey[:]...)
	message = append(message, newKey[:]...)
	signature := ed.Sign(privateKey, message)
	e := Entry{}
	e.ChainID = i.ChainID
	e.ExtIDs = [][]byte{[]byte("ReplaceKey"), oldKey[:], newKey[:], signature[:]}

	txID, err := CommitEntry(&e, ec)
	if err != nil {
		return "", err
	}
	_, err = RevealEntry(&e)
	if err != nil {
		return "", err
	}
	return txID, nil
}

func WriteAttribute(receiverChainID string, destinationChainID string, attributesJSON string, privateKey *[64]byte, signerChainID string, ec *ECAddress) (string, []byte, error) {
	message := []byte(receiverChainID + destinationChainID + attributesJSON)
	signature := ed.Sign(privateKey, message)
	publicKey := ed.GetPublicKey(privateKey)

	e := Entry{}
	e.ChainID = destinationChainID
	e.ExtIDs = [][]byte{[]byte("IdentityAttribute"), []byte(receiverChainID), signature[:], publicKey[:], []byte(signerChainID)}
	e.Content = []byte(attributesJSON)

	txID, err := CommitEntry(&e, ec)
	if err != nil {
		return "", nil, err
	}
	_, err = RevealEntry(&e)
	if err != nil {
		return "", nil, err
	}
	return txID, e.Hash(), nil
}

// EndorseIdentityAttribute signs a message using the provided private key saying that the signing identity acknowledges/agrees with
// the attribute entry located at entryHash
func EndorseIdentityAttribute(destinationChainID string, attributeEntryHash string, privateKey *[64]byte, signerChainID string, ec *ECAddress) (string, []byte, error) {
	message := []byte(destinationChainID + attributeEntryHash)
	signature := ed.Sign(privateKey, message)
	publicKey := ed.GetPublicKey(privateKey)

	e := Entry{}
	e.ChainID = destinationChainID
	e.ExtIDs = [][]byte{[]byte("IdentityAttributeEndorsement"), signature[:], publicKey[:], []byte(signerChainID)}
	e.Content = []byte(attributeEntryHash)

	txID, err := CommitEntry(&e, ec)
	if err != nil {
		return "", nil, err
	}
	_, err = RevealEntry(&e)
	if err != nil {
		return "", nil, err
	}
	return txID, e.Hash(), nil
}

// IsValidAttribute returns true if the EntryHash points to a correctly formatted attribute entry with a signature
// that was valid for its signer's identity at the time the attribute was published
func IsValidAttribute(entryHash string) (bool, error) {
	e, err := GetEntry(entryHash)
	if err != nil {
		return false, err
	}

	// Check ExtIDs for valid formatting, then process them
	if len(e.ExtIDs) < 5 || bytes.Compare(e.ExtIDs[0], []byte("IdentityAttribute")) != 0 {
		return false, nil
	}
	if len(e.ExtIDs[1]) != 64 || len(e.ExtIDs[2]) != 64 || len(e.ExtIDs[3]) != 32 || len(e.ExtIDs[4]) != 64 {
		return false, nil
	}
	receiverChainID := e.ExtIDs[1]
	var signature [64]byte
	copy(signature[:], e.ExtIDs[2])
	var pubKey [32]byte
	copy(pubKey[:], e.ExtIDs[3])
	signerChainID := string(e.ExtIDs[4])

	// Message that was signed = ReceiverChainID + DestinationChainID + AttributesJSON
	message := receiverChainID
	message = append(message, []byte(e.ChainID)...)
	message = append(message, e.Content...)
	if !ed.Verify(&pubKey, message, &signature) {
		return false, nil
	}

	// Check that public key was valid for the signer at the time of the attribute being published
	receipt, err := GetReceipt(entryHash)
	if err != nil {
		return false, err
	}
	dblock, err := GetDBlock(receipt.DirectoryBlockKeyMR)
	if err != nil {
		return false, err
	}
	validKeys, err := GetKeysAtHeight(signerChainID, dblock.Header.SequenceNumber)
	if err != nil {
		return false, err
	}
	for _, key := range validKeys {
		if bytes.Compare(pubKey[:], key.Pub[:]) == 0 {
			// Found provided key to be valid at time of publishing attribute
			return true, nil
		}
	}
	return false, nil
}

// IsValidEndorsement returns true if the EntryHash points to a correctly formatted endorsement entry with a signature
// that was valid for its signer's identity at the time the attribute was published
func IsValidEndorsement(entryHash string) (bool, error) {
	e, err := GetEntry(entryHash)
	if err != nil {
		return false, err
	}

	// Check ExtIDs for valid formatting, then process them
	if len(e.ExtIDs) < 4 || bytes.Compare(e.ExtIDs[0], []byte("IdentityAttributeEndorsement")) != 0 {
		return false, nil
	}
	if len(e.ExtIDs[1]) != 64 || len(e.ExtIDs[2]) != 32 || len(e.ExtIDs[3]) != 64 {
		return false, nil
	}
	var signature [64]byte
	copy(signature[:], e.ExtIDs[1])
	var pubKey [32]byte
	copy(pubKey[:], e.ExtIDs[2])
	signerChainID := string(e.ExtIDs[3])

	// Message that was signed = DestinationChainID + AttributeEntryHash
	message := []byte(e.ChainID)
	message = append(message, e.Content...)
	if !ed.Verify(&pubKey, message, &signature) {
		return false, nil
	}

	// Check that public key was valid for the signer at the time of the attribute being published
	receipt, err := GetReceipt(entryHash)
	if err != nil {
		return false, err
	}
	dblock, err := GetDBlock(receipt.DirectoryBlockKeyMR)
	if err != nil {
		return false, err
	}
	validKeys, err := GetKeysAtHeight(signerChainID, dblock.Header.SequenceNumber)
	if err != nil {
		return false, err
	}
	for _, key := range validKeys {
		if bytes.Compare(pubKey[:], key.Pub[:]) == 0 {
			// Found provided key to be valid at time of publishing attribute
			return true, nil
		}
	}
	return false, nil
}
