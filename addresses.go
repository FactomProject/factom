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
	"github.com/FactomProject/go-bip44"
)

type addressStringType byte

// Magic numbers representing values that AddressStringType() returns.
const (
	InvalidAddress addressStringType = iota
	FactoidPub
	FactoidSec
	ECPub
	ECSec
)

// Address structure length values
const (
	AddressLength  = 38
	PrefixLength   = 2
	ChecksumLength = 4
	BodyLength     = AddressLength - ChecksumLength
)

var (
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
	fcPubPrefix = []byte{0x5f, 0xb1}
	fcSecPrefix = []byte{0x64, 0x78}
)

// AddressStringType parses an address and returns an addressStringType
// representing the type of address. An address can be a public or private
// Factoid or Entry Credit address. See Constants.
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

// IsValidAddress takes a string and returns true if it is a valid address. A
// valid address has a correct prefix and ends with a valid checksum.
func IsValidAddress(s string) bool {
	p := base58.Decode(s)

	if len(p) != AddressLength {
		return false
	}

	prefix := p[:PrefixLength]
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
	body := p[:BodyLength]
	check := p[AddressLength-ChecksumLength:]
	if bytes.Equal(shad(body)[:ChecksumLength], check) {
		return true
	}

	return false
}

// ECAddress stores the bytes of the public and secret keys for an EC address.
type ECAddress struct {
	Pub *[ed.PublicKeySize]byte
	Sec *[ed.PrivateKeySize]byte
}

// NewECAddress returns an initialized but empty ECAddress struct.
func NewECAddress() *ECAddress {
	a := new(ECAddress)
	a.Pub = new([ed.PublicKeySize]byte)
	a.Sec = new([ed.PrivateKeySize]byte)
	return a
}

// UnmarshalBinary satisfies the BinaryUnmarshaler interface for ECAddress. It
// takes binary data and unmarshals it into an ECAddress.
func (a *ECAddress) UnmarshalBinary(data []byte) error {
	_, err := a.UnmarshalBinaryData(data)
	return err
}

// UnmarshalBinaryData takes binary data and unmarshals it into an ECAddress
// and returns a byte slice with the first 32 bytes of data removed.
func (a *ECAddress) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("secret key portion must be 32 bytes")
	}

	if a.Sec == nil {
		a.Sec = new([ed.PrivateKeySize]byte)
	}

	copy(a.Sec[:], data[:32])
	a.Pub = ed.GetPublicKey(a.Sec)

	return data[32:], nil
}

// MarshalBinary satisfies the BinaryMarshaler interface for ECAddress. It
// returns a byte slice containing the raw bytes of the secret key.
func (a *ECAddress) MarshalBinary() ([]byte, error) {
	return a.SecBytes()[:32], nil
}

// GetECAddress takes a private address string (Es...) and returns an
// ECAddress.
func GetECAddress(s string) (*ECAddress, error) {
	if !IsValidAddress(s) {
		return nil, fmt.Errorf("Invalid Address")
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:PrefixLength], ecSecPrefix) {
		return nil, fmt.Errorf("Invalid Entry Credit Private Address")
	}

	return MakeECAddress(p[PrefixLength:BodyLength])
}

// MakeECAddress takes a 32 byte secret key and returns an ECAddress
// representing that key.
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

// PubBytes returns the public key as a []byte.
func (a *ECAddress) PubBytes() []byte {
	return a.Pub[:]
}

// PubFixed returns the public key as a fixed size array.
func (a *ECAddress) PubFixed() *[ed.PublicKeySize]byte {
	return a.Pub
}

// PubString returns the string encoding of the public key.
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

// SecBytes returns the public key as a []byte.
func (a *ECAddress) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the public key as a fixed size array.
func (a *ECAddress) SecFixed() *[ed.PrivateKeySize]byte {
	return a.Sec
}

// SecString returns the string encoding of the secret key.
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

// Sign returns a signature of the msg with the ECAddress private key.
func (a *ECAddress) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(a.SecFixed(), msg)
}

// String satisfies the Stringer interface for ECAddress and returns the
// PubString.
func (a *ECAddress) String() string {
	return a.PubString()
}

// FactoidAddress stores
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

// UnmarshalBinary satisfies the BinaryMarshaler interface and populates a
// FactoidAddress from the binary data passed to it.
func (t *FactoidAddress) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

// UnmarshalBinaryData populates a FactoidAddress from the binary data passed
// to it.
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

// MarshalBinaryData returns a byte slice of the raw bytes of the secret key of
// the FactoidAddress.
func (t *FactoidAddress) MarshalBinary() ([]byte, error) {
	return t.SecBytes()[:32], nil
}

// GetFactoidAddress takes a private address string (Fs...) and returns a
// FactoidAddress.
func GetFactoidAddress(s string) (*FactoidAddress, error) {
	if !IsValidAddress(s) {
		return nil, fmt.Errorf("Invalid Address")
	}

	p := base58.Decode(s)

	if !bytes.Equal(p[:PrefixLength], fcSecPrefix) {
		return nil, fmt.Errorf("Invalid Factoid Private Address")
	}

	return MakeFactoidAddress(p[PrefixLength:BodyLength])
}

// MakeFactoidAddress makes a new factoid address using the sec bytes as the
// secret key.
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

// ParseAndValidateMnemonic takes a mnemonic strings and verifies that it is of
// the correct length and format. It returns the mnemonic with any whitespace
// trimmed and replaced with single spaces.
func ParseAndValidateMnemonic(mnemonic string) (string, error) {
	if l := len(strings.Fields(mnemonic)); l != 12 {
		return "", fmt.Errorf("Incorrect mnemonic length. Expecitng 12 words, found %d", l)
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

// MakeFactoidAddressFromKoinify takes the 12 word string used in the Koinify
// sale and returns a Factoid Address.
func MakeFactoidAddressFromKoinify(mnemonic string) (*FactoidAddress, error) {
	mnemonic, err := ParseAndValidateMnemonic(mnemonic)
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

func MakeBIP44FactoidAddress(mnemonic string, account, chain, address uint32) (*FactoidAddress, error) {
	mnemonic, err := ParseAndValidateMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeFactomFactoids, account, chain, address)
	if err != nil {
		return nil, err
	}

	return MakeFactoidAddress(child.Key)
}

func MakeBIP44ECAddress(mnemonic string, account, chain, address uint32) (*ECAddress, error) {
	mnemonic, err := ParseAndValidateMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeFactomEntryCredits, account, chain, address)
	if err != nil {
		return nil, err
	}

	return MakeECAddress(child.Key)
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
