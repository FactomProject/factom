package factom

import (
	"fmt"
	"testing"

	"encoding/json"
	ed "github.com/FactomProject/ed25519"
)

func TestGetIdentityChainID(t *testing.T) {
	name := []string{"John", "Jacob", "Jingleheimer-Schmidt"}
	observedChainID := GetIdentityChainID(name)
	expectedChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	if observedChainID != expectedChainID {
		t.Errorf("got: %s but expected: %s", observedChainID, expectedChainID)
	}
}

func TestNewIdentityChain(t *testing.T) {
	name := []string{"John", "Jacob", "Jingleheimer-Schmidt"}
	secretKeys := []string{
		"idsec2rChEHLz3SPQQx3syQtB11pHAmxyGjux5FntnS7xqTCieHxxTc",
		"idsec1xuUyeCCrJhsojf2wLAZqRxPzPFR8Gidd9DRRid1yGy8ncAJG3",
		"idsec2J3nNoqdiyboCBKDGauqN9Jb33dyFSqaJKZqTs6i5FmztsTn5f",
		"idsec1jztZ7dypqtwtPPWxybZFNpvvpUh6g8oog6Mnk2gGCm1pNBTgE",
	}
	var keys []*IdentityKey
	for _, v := range secretKeys {
		k, _ := GetIdentityKey(v)
		keys = append(keys, k)
	}

	newChain := NewIdentityChain(name, keys)
	expectedChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	t.Run("ChainID", func(t *testing.T) {
		if newChain.ChainID != expectedChainID {
			fmt.Println(newChain.ChainID)
			fmt.Println(expectedChainID)
			t.Fail()
		}
	})
	t.Run("Keys accessible from Content", func(t *testing.T) {
		var contentMap map[string]interface{}
		content := newChain.FirstEntry.Content
		err := json.Unmarshal(content, &contentMap)
		if err != nil {
			t.Errorf("Failed to unmarshal content")
		}
		for i, v := range contentMap["keys"].([]interface{}) {
			if keys[i].String() != v.(string) {
				t.Errorf("Keys not properly formatted")
			}
		}
	})

}

func TestNewIdentityKeyReplacementEntry(t *testing.T) {
	chainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	oldKey, _ := GetIdentityKey("idsec1jztZ7dypqtwtPPWxybZFNpvvpUh6g8oog6Mnk2gGCm1pNBTgE")
	newKey, _ := GetIdentityKey("idsec2J3nNoqdiyboCBKDGauqN9Jb33dyFSqaJKZqTs6i5FmztsTn5f")
	signerKey, _ := GetIdentityKey("idsec2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9h7s")

	observedEntry := NewIdentityKeyReplacementEntry(chainID, oldKey, newKey, signerKey)

	t.Run("ChainID", func(t *testing.T) {
		if observedEntry.ChainID != chainID {
			t.Fail()
		}
	})
	t.Run("ExtIDs", func(t *testing.T) {
		if len(observedEntry.ExtIDs) != 5 {
			t.Errorf("len(ExtIDs) != 5")
		}
		if string(observedEntry.ExtIDs[0]) != "ReplaceKey" {
			t.Errorf("ReplaceKey is not first ExtID")
		}
		if string(observedEntry.ExtIDs[1]) != oldKey.String() ||
			string(observedEntry.ExtIDs[2]) != newKey.String() ||
			string(observedEntry.ExtIDs[4]) != signerKey.String() {
			t.Errorf("Keys not formatted properly")
		}
	})
	t.Run("Signature", func(t *testing.T) {
		var observedSignature [64]byte
		copy(observedSignature[:], observedEntry.ExtIDs[3])
		message := []byte(oldKey.String() + newKey.String())
		if !ed.Verify(signerKey.Pub, message, &observedSignature) {
			t.Fail()
		}
	})
}

func TestNewIdentityAttributeEntry(t *testing.T) {
	receiverChainID := "5ef81cd345fd497a376ca5e5670ef10826d96e73c9f797b33ea46552a47834a3"
	destinationChainID := "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	signerChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	signerKey, err := GetIdentityKey("idsec2J3nNoqdiyboCBKDGauqN9Jb33dyFSqaJKZqTs6i5FmztsTn5f")
	if err != nil {
		t.Errorf("Failed to get identity key")
	}
	attributesJSON := `[{"key":"email","value":"abc@def.ghi"}]`

	observedEntry := NewIdentityAttributeEntry(receiverChainID, destinationChainID, attributesJSON, signerKey, signerChainID)

	if !IsValidAttribute(observedEntry) {
		t.Errorf("Improperly formatted attribute")
	}

	var attributes []IdentityAttribute
	if err = json.Unmarshal(observedEntry.Content, &attributes); err != nil {
		t.Errorf("Failed to unmarshal content: %v", err)
	}
	if attributes[0].Key != "email" {
		t.Errorf("Incorrect key")
	}
	if attributes[0].Value != "abc@def.ghi" {
		t.Errorf("Incorrect value")
	}
}

func TestNewIdentityAttributeEndorsementEntry(t *testing.T) {
	destinationChainID := "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	signerChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	signerKey, _ := GetIdentityKey("idsec2J3nNoqdiyboCBKDGauqN9Jb33dyFSqaJKZqTs6i5FmztsTn5f")
	entryHash := "52385948ea3ab6fd67b07664ac6a30ae5f6afa94427a547c142517beaa9054d0"

	observedEntry := NewIdentityAttributeEndorsementEntry(destinationChainID, entryHash, signerKey, signerChainID)

	if !IsValidEndorsement(observedEntry) {
		t.Errorf("Improperly formatted attribute endorsement")
	}
}
