package factom

import (
	"testing"
	"fmt"

	ed "github.com/FactomProject/ed25519"
	"encoding/json"
)

func TestGetIdentityChainID(t *testing.T) {
	name := []string{"John", "Jacob", "Jingleheimer-Schmidt"}
	observedChainID := GetIdentityChainID(name)
	expectedChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	if observedChainID != expectedChainID {
		fmt.Println(observedChainID)
		fmt.Println(expectedChainID)
		t.Fail()
	}
}

func TestNewIdentityChain(t *testing.T) {
	name := []string{"John", "Jacob", "Jingleheimer-Schmidt"}
	secretKeys := []string{
		"OQt+S8561HclsbNHEgPF6gE8ElOl+cQ/I9bqolAW2WE=",
		"IwlkC3xooRM2xq7N0m94InOYj9xcavQ36V3Ar3NUCW0=",
		"oZU/23DMo3UK4HctBUGkpdUqgt0UqeCF6BhlNdmcPwY=",
		"mwkcvN2a38Uv16FXlK01csWVjqXSvF7I6l+oBPmVCRQ=",
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
			fmt.Println("Failed to unmarshal content")
			t.Fail()
		}
		for i, v := range contentMap["keys"].([]interface{}) {
			if keys[i].String() != v.(string) {
				fmt.Println("Keys not properly formatted")
				t.Fail()
			}
		}
	})

}

func TestNewIdentityKeyReplacementEntry(t *testing.T) {
	chainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	oldKey, _ := GetIdentityKey("mwkcvN2a38Uv16FXlK01csWVjqXSvF7I6l+oBPmVCRQ=")
	newKey, _ := GetIdentityKey("TTqbfGahXE7MKJ1/kxv/HEGk0yAblehJ+tBs76goQAM=")
	signerKey, _ := GetIdentityKey("IwlkC3xooRM2xq7N0m94InOYj9xcavQ36V3Ar3NUCW0=")

	observedEntry := NewIdentityKeyReplacementEntry(chainID, oldKey, newKey, signerKey)

	t.Run("ChainID", func(t *testing.T) {
		if observedEntry.ChainID != chainID {
			t.Fail()
		}
	})
	t.Run("ExtIDs", func(t *testing.T) {
		if len(observedEntry.ExtIDs) != 5 {
			fmt.Println("len(ExtIDs) != 5")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[0]) != "ReplaceKey" {
			fmt.Println("ReplaceKey is not first ExtID")
			t.Fail()
		}
		if  string(observedEntry.ExtIDs[1]) != oldKey.String() ||
			string(observedEntry.ExtIDs[2]) != newKey.String() ||
			string(observedEntry.ExtIDs[4]) != signerKey.String() {
			fmt.Println("Keys not formatted properly")
			t.Fail()
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
	signerKey, _ := GetIdentityKey("IwlkC3xooRM2xq7N0m94InOYj9xcavQ36V3Ar3NUCW0=")
	attributes := `[{"key":"email","value":"abc@def.ghi"}]`

	observedEntry := NewIdentityAttributeEntry(receiverChainID, destinationChainID, attributes, signerKey, signerChainID)

	t.Run("ChainID", func(t *testing.T) {
		if observedEntry.ChainID != destinationChainID {
			fmt.Println("Incorrect Destination ChainID")
			t.Fail()
		}
	})
	t.Run("ExtIDs", func(t *testing.T) {
		if len(observedEntry.ExtIDs) != 5 {
			fmt.Println("len(ExtIDs) != 5")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[0]) != "IdentityAttribute" {
			fmt.Println("IdentityAttribute is not first ExtID")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[1]) != receiverChainID {
			fmt.Println("Receiver ChainID is not ExtID[1]")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[4]) != signerChainID {
			fmt.Println("Signer ChainID is not ExtID[4]")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[3]) != signerKey.String() {
			fmt.Println("Signer key not properly formatted or is not ExtID[3]")
			t.Fail()
		}
	})
	t.Run("Attributes accessible from Content", func(t *testing.T) {
		var attributes []IdentityAttribute
		err := json.Unmarshal(observedEntry.Content, &attributes)
		if err != nil {
			fmt.Println("Failed to unmarshal content")
		}
		if attributes[0].Key != "email" {
			fmt.Println("Incorrect key")
			t.Fail()
		}
		if attributes[0].Value != "abc@def.ghi" {
			fmt.Println("Incorrect value")
			t.Fail()
		}
	})
	t.Run("Signature", func(t *testing.T) {
		var observedSignature [64]byte
		copy(observedSignature[:], observedEntry.ExtIDs[2])
		message := []byte(receiverChainID + destinationChainID + attributes)
		if !ed.Verify(signerKey.Pub, message, &observedSignature) {
			t.Fail()
		}
	})
}

func TestNewIdentityAttributeEndorsementEntry(t *testing.T) {
	destinationChainID := "5a402200c5cf278e47905ce52d7d64529a0291829a7bd230072c5468be709069"
	signerChainID := "e0cf1713b492e09e783d5d9f4fc6e2c71b5bdc9af4806a7937a5e935819717e9"
	signerKey, _ := GetIdentityKey("IwlkC3xooRM2xq7N0m94InOYj9xcavQ36V3Ar3NUCW0=")
	entryHash := "52385948ea3ab6fd67b07664ac6a30ae5f6afa94427a547c142517beaa9054d0"

	observedEntry := NewIdentityAttributeEndorsementEntry(destinationChainID, entryHash, signerKey, signerChainID)

	t.Run("ChainID", func(t *testing.T) {
		if observedEntry.ChainID != destinationChainID {
			fmt.Println("Incorrect Destination ChainID")
			t.Fail()
		}
	})
	t.Run("ExtIDs", func(t *testing.T) {
		if len(observedEntry.ExtIDs) != 4 {
			fmt.Println("len(ExtIDs) != 4")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[0]) != "IdentityAttributeEndorsement" {
			fmt.Println("IdentityAttributeEndorsement is not first ExtID")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[3]) != signerChainID {
			fmt.Println("Signer ChainID is not ExtID[3]")
			t.Fail()
		}
		if string(observedEntry.ExtIDs[2]) != signerKey.String() {
			fmt.Println("Signer key not properly formatted or is not ExtID[2]")
			t.Fail()
		}
	})
	t.Run("Signature", func(t *testing.T) {
		var observedSignature [64]byte
		copy(observedSignature[:], observedEntry.ExtIDs[1])
		message := []byte(destinationChainID + entryHash)
		if !ed.Verify(signerKey.Pub, message, &observedSignature) {
			t.Fail()
		}
	})
}