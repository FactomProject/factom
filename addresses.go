// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"errors"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/go-bip32"
	"github.com/FactomProject/go-bip39"
	"github.com/FactomProject/go-bip44"
)

// Common Address errors
var (
	ErrInvalidAddress    = errors.New("invalid address")
	ErrInvalidFactoidSec = errors.New("invalid Factoid secret address")
	ErrInvalidECSec      = errors.New("invalid Entry Credit secret address")
	ErrSecKeyLength      = errors.New("secret key portion must be 32 bytes")
	ErrMnemonicLength    = errors.New("mnemonic must be 12 words")
)

type addressStringType byte

const (
	InvalidAddress addressStringType = iota
	FactoidPub
	FactoidSec
	ECPub
	ECSec
)

const (
	AddressLength  = 38
	PrefixLength   = 2
	ChecksumLength = 4
	BodyLength     = AddressLength - ChecksumLength
)

var (
	fcPubPrefix = []byte{0x5f, 0xb1}
	fcSecPrefix = []byte{0x64, 0x78}
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
)

func AddressStringType(s string) addressStringType {
	p := base58.Decode(s)

	if len(p) != AddressLength {
		return InvalidAddress
	}

	// verify the address checksum
	body := p[:BodyLength]
	check := p[AddressLength-ChecksumLength:]
	if !bytes.Equal(shad(body)[:ChecksumLength], check) {
		return InvalidAddress
	}

	prefix := p[:PrefixLength]
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
	if AddressStringType(s) != InvalidAddress {
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
		return nil, ErrSecKeyLength
	}

	if a.Sec == nil {
		a.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(a.Sec[:], data[:32])
	a.Pub = ed.GetPublicKey(a.Sec)

	return data[32:], nil
}

func (a *ECAddress) MarshalBinary() ([]byte, error) {
	return a.SecBytes()[:32], nil
}

// GetECAddress takes a private address string (Es...) and returns an ECAddress.
func GetECAddress(s string) (*ECAddress, error) {
	if !IsValidAddress(s) {
		return nil, ErrInvalidAddress
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:PrefixLength], ecSecPrefix) {
		return nil, ErrInvalidECSec
	}

	return MakeECAddress(p[PrefixLength:BodyLength])
}

func MakeECAddress(sec []byte) (*ECAddress, error) {
	if len(sec) != 32 {
		return nil, ErrSecKeyLength
	}

	a := NewECAddress()

	err := a.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func MakeBIP44ECAddress(mnemonic string, account, chain, address uint32) (*ECAddress, error) {
	mnemonic, err := ParseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeFactomEntryCredits, account, chain, address)
	if err != nil {
		return nil, err
	}

	return MakeECAddress(child.Key)
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
	check := shad(buf.Bytes())[:ChecksumLength]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

// SecBytes returns the []byte representation of the secret key
func (a *ECAddress) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the fixed size secret key
func (a *ECAddress) SecFixed() *[ed.PrivateKeySize]byte {
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
	check := shad(buf.Bytes())[:ChecksumLength]
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

// GetFactoidAddress takes a private address string (Fs...) and returns a
// FactoidAddress.
func GetFactoidAddress(s string) (*FactoidAddress, error) {
	if !IsValidAddress(s) {
		return nil, ErrInvalidAddress
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:PrefixLength], fcSecPrefix) {
		return nil, ErrInvalidFactoidSec
	}

	return MakeFactoidAddress(p[PrefixLength:BodyLength])
}

func MakeFactoidAddress(sec []byte) (*FactoidAddress, error) {
	if len(sec) != 32 {
		return nil, ErrSecKeyLength
	}

	a := NewFactoidAddress()
	err := a.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func MakeBIP44FactoidAddress(mnemonic string, account, chain, address uint32) (*FactoidAddress, error) {
	mnemonic, err := ParseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeFactomFactoids, account, chain, address)
	if err != nil {
		return nil, err
	}

	return MakeFactoidAddress(child.Key)
}

// MakeFactoidAddressFromKoinify takes the 12 word string used in the Koinify
// sale and returns a Factoid Address.
func MakeFactoidAddressFromKoinify(mnemonic string) (*FactoidAddress, error) {
	mnemonic, err := ParseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

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

func (a *FactoidAddress) UnmarshalBinary(data []byte) error {
	_, err := a.UnmarshalBinaryData(data)
	return err
}

func (a *FactoidAddress) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, ErrSecKeyLength
	}

	if a.Sec == nil {
		a.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(a.Sec[:], data[:32])
	r := NewRCD1()
	r.Pub = ed.GetPublicKey(a.Sec)
	a.RCD = r

	return data[32:], nil
}

func (a *FactoidAddress) MarshalBinary() ([]byte, error) {
	return a.SecBytes()[:32], nil
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

func (a *FactoidAddress) SecFixed() *[ed.PrivateKeySize]byte {
	return a.Sec
}

func (a *FactoidAddress) SecString() string {
	buf := new(bytes.Buffer)

	// Factoid address prefix
	buf.Write(fcSecPrefix)

	// Secret key
	buf.Write(a.SecBytes()[:32])

	// Checksum
	check := shad(buf.Bytes())[:ChecksumLength]
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
	check := shad(buf.Bytes())[:ChecksumLength]
	buf.Write(check)

	return base58.Encode(buf.Bytes())
}

func ParseMnemonic(mnemonic string) (string, error) {
	if l := len(strings.Fields(mnemonic)); l != 12 {
		return "", ErrMnemonicLength
	}

	mnemonic = strings.ToLower(strings.TrimSpace(mnemonic))

	split := strings.Split(mnemonic, " ")
	for i := len(split) - 1; i >= 0; i-- {
		if split[i] == "" {
			split = append(split[:i], split[i+1:]...)
		}
	}
	mnemonic = strings.Join(split, " ")

	_, err := bip39.MnemonicToByteArray(mnemonic)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
