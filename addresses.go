// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/go-bip32"
	"github.com/FactomProject/go-bip39"
)

type addressStringType byte

const (
	InvalidAddress addressStringType = iota
	FactoidPub
	FactoidSec
	ECPub
	ECSec
)

var (
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
	fcPubPrefix = []byte{0x5f, 0xb1}
	fcSecPrefix = []byte{0x64, 0x78}
)

func AddressStringType(s string) addressStringType {
	p := base58.Decode(s)

	if len(p) != 38 {
		return InvalidAddress
	}

	// verify the address checksum
	body := p[:len(p)-4]
	check := p[len(p)-4:]
	if !bytes.Equal(shad(body)[:4], check) {
		return InvalidAddress
	}

	prefix := p[:2]
	switch {
	case bytes.Equal(prefix, ecPubPrefix):
		return ECPub
	case bytes.Equal(prefix, ecSecPrefix):
		return ECSec
	case bytes.Equal(prefix, fcPubPrefix):
		return FactoidPub
	case bytes.Equal(prefix, fcSecPrefix):
		return FactoidSec
	default:
		return InvalidAddress
	}
}

func IsValidAddress(s string) bool {
	p := base58.Decode(s)

	if len(p) != 38 {
		return false
	}

	prefix := p[:2]
	switch {
	case bytes.Equal(prefix, ecPubPrefix):
		break
	case bytes.Equal(prefix, ecSecPrefix):
		break
	case bytes.Equal(prefix, fcPubPrefix):
		break
	case bytes.Equal(prefix, fcSecPrefix):
		break
	default:
		return false
	}

	// verify the address checksum
	body := p[:len(p)-4]
	check := p[len(p)-4:]
	if bytes.Equal(shad(body)[:4], check) {
		return true
	}

	return false
}

type ECAddress struct {
	Pub *[ed.PublicKeySize]byte
	Sec *[ed.PrivateKeySize]byte
}

func NewECAddress() *ECAddress {
	a := new(ECAddress)
	a.Pub = new([ed.PublicKeySize]byte)
	a.Sec = new([ed.PrivateKeySize]byte)
	return a
}

func (a *ECAddress) UnmarshalBinary(data []byte) error {
	_, err := a.UnmarshalBinaryData(data)
	return err
}

func (a *ECAddress) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	if a.Sec == nil {
		a.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(t.Sec[:], data[:32])
	a.Pub = ed.GetPublicKey(a.Sec)

	return data[32:], nil
}

func (a *ECAddress) MarshalBinary() ([]byte, error) {
	return a.SecBytes(), nil
}

// GetECAddress takes a private address string (Es...) and returns an ECAddress.
func GetECAddress(s string) (*ECAddress, error) {
	if !IsValidAddress(s) {
		return nil, fmt.Errorf("Invalid Address")
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:2], ecSecPrefix) {
		return nil, fmt.Errorf("Invalid Entry Credit Private Address")
	}

	return MakeECAddress(p[2:34])
}

func MakeECAddress(sec []byte) (*ECAddress, error) {
	if len(sec) != 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	a := NewECAddress()

	err := a.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// PubBytes returns the []byte representation of the public key
func (a *ECAddress) PubBytes() []byte {
	return a.Pub[:]
}

// PubFixed returns the fixed size public key
func (a *ECAddress) PubFixed() *[ed.PublicKeySize]byte {
	return a.Pub
}

// PubString returns the string encoding of the public key
func (a *ECAddress) PubString() string {
	buf := new(bytes.Buffer)

	// EC address prefix
	buf.Write(ecPubPrefix)

	// Public key
	buf.Write(a.PubBytes())

	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

// SecBytes returns the []byte representation of the secret key
func (a *ECAddress) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the fixed size secret key
func (a *ECAddress) SecFixed() *[64]byte {
	return a.Sec
}

// SecString returns the string encoding of the secret key
func (a *ECAddress) SecString() string {
	buf := new(bytes.Buffer)

	// EC address prefix
	buf.Write(ecSecPrefix)

	// Secret key
	buf.Write(a.SecBytes()[:32])

	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

// Sign the message with the ECAddress private key
func (a *ECAddress) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(a.SecFixed(), msg)
}

func (a *ECAddress) String() string {
	return a.PubString()
}

type FactoidAddress struct {
	RCD RCD
	Sec *[ed.PrivateKeySize]byte
}

func NewFactoidAddress() *FactoidAddress {
	a := new(FactoidAddress)
	r := NewRCD1()
	r.Pub = new([ed.PublicKeySize]byte)
	a.RCD = r
	a.Sec = new([ed.PrivateKeySize]byte)
	return a
}

func (t *FactoidAddress) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (t *FactoidAddress) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	if t.Sec == nil {
		t.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(t.Sec[:], data[:32])
	r := NewRCD1()
	r.Pub = ed.GetPublicKey(t.Sec)
	t.RCD = r

	return data[32:], nil
}

func (t *FactoidAddress) MarshalBinary() ([]byte, error) {
	return t.SecBytes(), nil
}

// GetFactoidAddress takes a private address string (Fs...) and returns a
// FactoidAddress.
func GetFactoidAddress(s string) (*FactoidAddress, error) {
	if !IsValidAddress(s) {
		return nil, fmt.Errorf("Invalid Address")
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:2], fcSecPrefix) {
		return nil, fmt.Errorf("Invalid Factoid Private Address")
	}

	return MakeFactoidAddress(p[2:34])
}

func MakeFactoidAddress(sec []byte) (*FactoidAddress, error) {
	if len(sec) != 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	a := NewFactoidAddress()
	err := a.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// MakeFactoidAddressFromMnemonic takes the 12 word string used in the Koinify
// sale and returns a Factoid Address.
func MakeFactoidAddressFromMnemonic(mnemonic string) (*FactoidAddress, error) {
	l := len(strings.Fields(mnemonic))
	if l < 12 {
		return nil, fmt.Errorf("Not enough words in mnemonic. Expecitng 12, found %d", l)
	}
	if l > 12 {
		return nil, fmt.Errorf("Too many words in mnemonic. Expecitng 12, found %d", l)
	}

	mnemonic = strings.ToLower(strings.TrimSpace(mnemonic))
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	child, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 7)
	if err != nil {
		return nil, err
	}

	return MakeFactoidAddress(child.Key)
}

func (a *FactoidAddress) RCDHash() []byte {
	return a.RCD.Hash()
}

func (a *FactoidAddress) RCDType() uint8 {
	return a.RCD.Type()
}

func (a *FactoidAddress) PubBytes() []byte {
	return a.RCD.(*RCD1).PubBytes()
}

func (a *FactoidAddress) SecBytes() []byte {
	return a.Sec[:]
}

func (a *FactoidAddress) SecFixed() *[64]byte {
	return a.Sec
}

func (a *FactoidAddress) SecString() string {
	buf := new(bytes.Buffer)

	// Factoid address prefix
	buf.Write(fcSecPrefix)

	// Secret key
	buf.Write(a.SecBytes()[:32])

	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

func (a *FactoidAddress) String() string {
	buf := new(bytes.Buffer)

	// FC address prefix
	buf.Write(fcPubPrefix)

	// RCD Hash
	buf.Write(a.RCDHash())

	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}
