package factom_test

import (
	"testing"

	"bytes"
	"crypto/rand"

	. "github.com/FactomProject/factom"

	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/go-bip32"
)

func TestMarshalIdendityKey(t *testing.T) {
	for i := 0; i < 100; i++ {
		sec := make([]byte, 32)
		_, err := rand.Read(sec)
		if err != nil {
			t.Error(err)
		}

		k1, err := MakeIdentityKey(sec)
		if err != nil {
			t.Error(err)
		}

		data1, err := k1.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		k2 := new(ECAddress)
		data2, err := k2.UnmarshalBinaryData(data1)
		if err != nil {
			t.Error(err)
		}

		if len(data2) != 0 {
			t.Errorf("UnmarshalBinary left %d bytes remaining", len(data2))
		}

		if bytes.Compare(k1.SecBytes(), k2.SecBytes()) != 0 {
			t.Errorf("Unmarshaled object has different secret.")
		}

		if bytes.Compare(k1.PubBytes(), k2.PubBytes()) != 0 {
			t.Errorf("Unmarshaled object has different public.")
		}
	}
}

func TestNewIdentityKey(t *testing.T) {
	key := NewIdentityKey()
	pub := "idpub1koTMq9h7FRCAdgZmDhjW85FBUsDJ8n1MBz94UaWf61JvLL1aa"
	if key.PubString() != pub {
		t.Errorf("new pubkey %s did not match %s", key.PubString(), pub)
	}
}

func TestGetIdentityKey(t *testing.T) {
	pub := "idpub1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tByb"
	sec := "idsec2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9h7s"

	k, err := GetIdentityKey(sec)
	if err != nil {
		t.Error(err)
	}

	if k.PubString() != pub {
		t.Errorf("%s did not match %s", k.PubString(), pub)
	}

	msg := []byte("Hello Factom!")
	sig := k.Sign(msg)
	if !ed.Verify(k.PubFixed(), msg, sig) {
		t.Errorf("Key signature did not match")
	}
}

func TestMakeBIP44IdentityKey(t *testing.T) {
	m := "yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow yellow"
	pub := "idpub1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tByb"
	sec := "idsec2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9h7s"

	id, err := MakeBIP44IdentityKey(m, bip32.FirstHardenedChild, 0, 0)
	if err != nil {
		t.Error(err)
	}

	if id.String() != pub {
		t.Errorf("incorrect public key from 12 words: got %s expecting %s", id.String(), pub)
	}
	if id.SecString() != sec {
		t.Errorf("incorrect secret key from 12 words: got %s expecting %s", id.SecString(), sec)
	}
}

func TestIsValidIdentityKey(t *testing.T) {
	pub := "idpub1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tByb"
	sec := "idsec2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9h7s"
	badEmpty := ""
	badLen := "idpub1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tBybd"
	badPrePub := "idpXb1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tByb"
	badPreSec := "idsXc2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9h7s"
	badCheckPub := "idpub1p4YkMzskVrtbK45nBHaikGda9w5SMvKvVsQtgVUfLK5Y8tBby"
	badCheckSec := "idsec2wH72BNR9QZhTMGDbxwLWGrghZQexZvLTros2wCekkc62N9hs7"

	if !IsValidIdentityKey(pub) {
		t.Errorf("%s was not considered valid", pub)
	}
	if !IsValidIdentityKey(sec) {
		t.Errorf("%s was not considered valid", sec)
	}

	if IsValidIdentityKey(badEmpty) {
		t.Errorf("%s was considered valid", badEmpty)
	}
	if IsValidIdentityKey(badLen) {
		t.Errorf("%s was considered valid", badLen)
	}
	if IsValidIdentityKey(badPrePub) {
		t.Errorf("%s was considered valid", badPrePub)
	}
	if IsValidIdentityKey(badPreSec) {
		t.Errorf("%s was considered valid", badPreSec)
	}
	if IsValidIdentityKey(badCheckPub) {
		t.Errorf("%s was considered valid", badCheckPub)
	}
	if IsValidIdentityKey(badCheckSec) {
		t.Errorf("%s was considered valid", badCheckSec)
	}
}
