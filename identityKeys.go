package factom

import (
	"bytes"
	"fmt"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/go-bip44"
)

type identityKeyStringType byte

const (
	InvalidIdentityKey identityKeyStringType = iota
	IDPub
	IDSec
)

const (
	IDKeyLength       = 41
	IDKeyPrefixLength = 5
	IDKeyBodyLength   = IDKeyLength - ChecksumLength
)

var (
	idPubPrefix = []byte{0x03, 0x45, 0xef, 0x9d, 0xe0}
	idSecPrefix = []byte{0x03, 0x45, 0xf3, 0xd0, 0xd6}
)

func IdentityKeyStringType(s string) identityKeyStringType {
	p := base58.Decode(s)

	if len(p) != IDKeyLength {
		return InvalidIdentityKey
	}

	// verify the address checksum
	body := p[:IDKeyBodyLength]
	check := p[IDKeyLength-ChecksumLength:]
	if !bytes.Equal(shad(body)[:ChecksumLength], check) {
		return InvalidIdentityKey
	}

	prefix := p[:IDKeyPrefixLength]
	switch {
	case bytes.Equal(prefix, idPubPrefix):
		return IDPub
	case bytes.Equal(prefix, idSecPrefix):
		return IDSec
	default:
		return InvalidIdentityKey
	}
}

func IsValidIdentityKey(s string) bool {
	p := base58.Decode(s)

	if len(p) != IDKeyLength {
		return false
	}

	prefix := p[:IDKeyPrefixLength]
	switch {
	case bytes.Equal(prefix, idPubPrefix):
		break
	case bytes.Equal(prefix, idSecPrefix):
		break
	default:
		return false
	}

	// verify the address checksum
	body := p[:IDKeyBodyLength]
	check := p[IDKeyLength-ChecksumLength:]
	if bytes.Equal(shad(body)[:ChecksumLength], check) {
		return true
	}

	return false
}

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
	if !IsValidIdentityKey(s) {
		return nil, fmt.Errorf("invalid Identity Private Key")
	}
	p := base58.Decode(s)

	if !bytes.Equal(p[:IDKeyPrefixLength], idSecPrefix) {
		return nil, fmt.Errorf("invalid Identity Private Key")
	}

	return MakeIdentityKey(p[IDKeyPrefixLength:IDKeyBodyLength])
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

func MakeBIP44IdentityKey(mnemonic string, account, chain, address uint32) (*IdentityKey, error) {
	mnemonic, err := ParseAndValidateMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeFactomIdentitiy, account, chain, address)
	if err != nil {
		return nil, err
	}

	return MakeIdentityKey(child.Key)
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
	buf := new(bytes.Buffer)
	buf.Write(idPubPrefix)
	buf.Write(k.PubBytes())

	check := shad(buf.Bytes())[:ChecksumLength]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
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
	buf := new(bytes.Buffer)
	buf.Write(idSecPrefix)
	buf.Write(k.SecBytes()[:32])

	check := shad(buf.Bytes())[:ChecksumLength]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

// Sign the message with the Identity private key
func (k *IdentityKey) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(k.SecFixed(), msg)
}

func (k *IdentityKey) String() string {
	return k.PubString()
}
