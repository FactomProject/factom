package factom

import (
	"fmt"

	"encoding/base64"

	ed "github.com/FactomProject/ed25519"
)

type IdentityKey struct {
	Pub *[ed.PublicKeySize]byte
	Sec *[ed.PrivateKeySize]byte
}

func NewIdentityKey() *IdentityKey {
	k := new(IdentityKey)
	k.Pub = new([ed.PublicKeySize]byte)
	k.Sec = new([ed.PrivateKeySize]byte)
	return k
}

func (k *IdentityKey) UnmarshalBinary(data []byte) error {
	_, err := k.UnmarshalBinaryData(data)
	return err
}

func (k *IdentityKey) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	if k.Sec == nil {
		k.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(k.Sec[:], data[:32])
	k.Pub = ed.GetPublicKey(k.Sec)

	return data[32:], nil
}

func (k *IdentityKey) MarshalBinary() ([]byte, error) {
	return k.SecBytes()[:32], nil
}

// GetIdentityKey takes a private key string and returns an IdentityKey.
func GetIdentityKey(s string) (*IdentityKey, error) {
	sec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode identity private key string")
	}
	if len(sec) != 32 {
		return nil, fmt.Errorf("incorrect length for identity private key string")
	}

	return MakeIdentityKey(sec)
}

func MakeIdentityKey(sec []byte) (*IdentityKey, error) {
	if len(sec) != 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	k := NewIdentityKey()

	err := k.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// PubBytes returns the []byte representation of the public key
func (k *IdentityKey) PubBytes() []byte {
	return k.Pub[:]
}

// PubFixed returns the fixed size public key
func (k *IdentityKey) PubFixed() *[ed.PublicKeySize]byte {
	return k.Pub
}

// PubString returns the string encoding of the public key
func (k *IdentityKey) PubString() string {
	return base64.StdEncoding.EncodeToString(k.PubBytes())
}

// SecBytes returns the []byte representation of the secret key
func (k *IdentityKey) SecBytes() []byte {
	return k.Sec[:]
}

// SecFixed returns the fixed size secret key
func (k *IdentityKey) SecFixed() *[ed.PrivateKeySize]byte {
	return k.Sec
}

// SecString returns the string encoding of the secret key
func (k *IdentityKey) SecString() string {
	return base64.StdEncoding.EncodeToString(k.SecBytes()[:32])
}

// Sign the message with the Identity private key
func (k *IdentityKey) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(k.SecFixed(), msg)
}

func (k *IdentityKey) String() string {
	return k.PubString()
}
