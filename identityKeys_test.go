package factom

import (
	"crypto/rand"
	"testing"

	ed "github.com/FactomProject/ed25519"
	"bytes"
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
	pub := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	if key.PubString() != pub {
		t.Errorf("new pubkey %s did not match %s", key.PubString(), pub)
	}
}

func TestGetIdentityKey(t *testing.T) {
	pub := "l7He0I0ziouQ3ffwRKfiDAI+82sZ51XQF9/W0Smhnjs="
	sec := "aokN4TYcmHBxP4WTYCan0ymYQtqLoLeKbnNWivylD8g="

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
