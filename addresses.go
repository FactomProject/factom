// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
	"github.com/FactomProject/go-bip32"
	"github.com/FactomProject/go-bip39"
	"github.com/FactomProject/go-bip44"
	"github.com/ethereum/go-ethereum/crypto"
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
	EthSec // 0x[32byte hex]
	EthFA  // eFA..
)

const (
	AddressLength  = 38
	PrefixLength   = 2
	ChecksumLength = 4
	BodyLength     = AddressLength - ChecksumLength

	// In hex characters, not bytes
	EthSecretLength = 64
	EthSecretPrefix = 2
)

var (
	fcPubPrefix = []byte{0x5f, 0xb1}
	fcSecPrefix = []byte{0x64, 0x78}
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
)

// AddressStringType determine the type of address from the given string.
// AddressStringType must return one of the defined address types;
// InvalidAddress, FactoidPub, FactoidSec, ECPub, or ECSec.
func AddressStringType(s string) addressStringType {
	if has0xPrefix(s) {
		_, err := hex.DecodeString(s[2:])
		if err != nil {
			return InvalidAddress
		}

		// Prefix + Secret length check means secret key
		if len(s) == EthSecretPrefix+EthSecretLength {
			return EthSec
		}

		return InvalidAddress
	}

	if hasEFAPrefix(s) {
		ty := AddressStringType(s[1:]) // Chop the e/E
		if ty == FactoidPub {
			return EthFA
		}
		return InvalidAddress
	}

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

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// has0xPrefix validates str begins with 'eFA' or 'EFA'.
func hasEFAPrefix(str string) bool {
	return len(str) >= 3 &&
		(str[0] == 'e' || str[0] == 'E') &&
		(str[1] == 'F') && (str[2] == 'A')
}

// IsValidAddress checks that a string is a valid address of one of the defined
// address types.
//
// For an address to be valid it must be the correct length, it must begin with
// one of the defined address prefixes, and the address checksum must match the
// address body.
func IsValidAddress(s string) bool {
	if AddressStringType(s) != InvalidAddress {
		return true
	}
	return false
}

// ECAddress is an Entry Credit public/secret key pair.
type ECAddress struct {
	Pub *[ed.PublicKeySize]byte
	Sec *[ed.PrivateKeySize]byte
}

// NewECAddress creates a blank public/secret key pair for an Entry Credit
// Address.
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

// UnmarshalBinaryData reads an ECAddress from a byte stream and returns the
// remainder of the byte stream.
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

// GetECAddress creates an Entry Credit Address public/secret key pair from a
// secret Entry Credit Address string i.e. Es...
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

// MakeECAddress creates an Entry Credit Address public/secret key pair from a
// secret key []byte.
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

// MakeBIP44ECAddress generates an Entry Credit Address from a 12 word mnemonic,
// an account index, a chain index, and an address index, according to the bip44
// standard for multicoin wallets.
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

// PubBytes returns the []byte representation of the public key.
func (a *ECAddress) PubBytes() []byte {
	return a.Pub[:]
}

// PubFixed returns the fixed size public key ([32]byte).
func (a *ECAddress) PubFixed() *[ed.PublicKeySize]byte {
	return a.Pub
}

// PubString returns the string encoding of the public key i.e. EC...
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

// SecBytes returns the []byte representation of the secret key.
func (a *ECAddress) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the fixed size secret key ([64]byte).
func (a *ECAddress) SecFixed() *[ed.PrivateKeySize]byte {
	return a.Sec
}

// SecString returns the string encoding of the secret key i.e. Es...
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

// Sign the message with the ECAddress secret key.
func (a *ECAddress) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(a.SecFixed(), msg)
}

func (a *ECAddress) String() string {
	return a.PubString()
}

// FactoidAddress is a Factoid Redeem Condition Datastructure (a type 1 RCD is
// just the public key) and a corresponding secret key.
type FactoidAddress struct {
	RCD RCD
	Sec *[ed.PrivateKeySize]byte
}

// NewFactoidAddress creates a blank rcd/secret key pair for a Factoid Address.
func NewFactoidAddress() *FactoidAddress {
	a := new(FactoidAddress)
	r := NewRCD1()
	r.Pub = new([ed.PublicKeySize]byte)
	a.RCD = r
	a.Sec = new([ed.PrivateKeySize]byte)
	return a
}

// GetFactoidAddress creates a Factoid Address rcd/secret key pair from a secret
// Factoid Address string i.e. Fs...
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

// MakeFactoidAddress creates a Factoid Address rcd/secret key pair from a
// secret key []byte.
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

// MakeBIP44FactoidAddress generates a Factoid Address from a 12 word mnemonic,
// an account index, a chain index, and an address index, according to the bip44
// standard for multicoin wallets.
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

// RCDHash returns the Hash of the Redeem Condition Datastructure from a Factoid
// Address.
func (a *FactoidAddress) RCDHash() []byte {
	return a.RCD.Hash()
}

// RCDType returns the Redeem Condition Datastructure type used by the Factoid
// Address.
func (a *FactoidAddress) RCDType() uint8 {
	return a.RCD.Type()
}

// PubBytes returns the []byte representation of the Redeem Condition
// Datastructure.
func (a *FactoidAddress) PubBytes() []byte {
	return a.RCD.(*RCD1).PubBytes()
}

// SecBytes returns the []byte representation of the secret key.
func (a *FactoidAddress) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the fixed size secret key ([64]byte).
func (a *FactoidAddress) SecFixed() *[ed.PrivateKeySize]byte {
	return a.Sec
}

// SecString returns the string encoding of the secret key i.e. Es...
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

// EthSecret is a Factoid Redeem Condition Datastructure (a type 0x0e RCD is
// just the public key) and a corresponding secret key.
type EthSecret struct {
	RCD RCD
	Sec *[32]byte
}

// NewFactoidAddress creates a blank rcd/secret key pair for a Factoid Address.
func NewEthSecret() *EthSecret {
	a := new(EthSecret)
	r := NewRCDe()
	r.Pub = new([64]byte)
	a.RCD = r
	a.Sec = new([32]byte)
	return a
}

func (a EthSecret) PrivateKey() *ecdsa.PrivateKey {
	secret, err := crypto.ToECDSA(a.Sec[:])
	if err != nil {
		return nil
	}
	return secret
}

func (a EthSecret) PublicKey() ecdsa.PublicKey {
	secret := a.PrivateKey()
	if secret == nil {
		return ecdsa.PublicKey{}
	}
	return secret.PublicKey
}

// GetEthSecret creates a Factoid Address rcd/secret key pair from a secret
// Factoid Address string i.e. 0x1234...
func GetEthSecret(s string) (*EthSecret, error) {
	if AddressStringType(s) != EthSec {
		return nil, ErrInvalidAddress
	}

	sec, err := hex.DecodeString(s[2:])
	if err != nil { // Should never hit this
		return nil, ErrInvalidAddress
	}

	return MakeEthSecret(sec)
}

// MakeEthSecret creates a EthSecret rcd/secret key pair from a
// secret key []byte.
func MakeEthSecret(sec []byte) (*EthSecret, error) {
	if len(sec) != 32 {
		return nil, ErrSecKeyLength
	}

	a := NewEthSecret()
	err := a.UnmarshalBinary(sec)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// MakeBIP44EthSecret generates an EthSecret from a 12 word mnemonic,
// an account index, a chain index, and an address index, according to the bip44
// standard for multicoin wallets.
func MakeBIP44EthSecret(mnemonic string, account, chain, address uint32) (*EthSecret, error) {
	mnemonic, err := ParseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	child, err := bip44.NewKeyFromMnemonic(mnemonic, bip44.TypeEther, account, chain, address)
	if err != nil {
		return nil, err
	}

	// TODO: Verify this against eth derivation
	return MakeEthSecret(child.Key)
}

func (a *EthSecret) UnmarshalBinary(data []byte) error {
	_, err := a.UnmarshalBinaryData(data)
	return err
}

func (a *EthSecret) UnmarshalBinaryData(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, ErrSecKeyLength
	}

	if a.Sec == nil {
		a.Sec = new([32]byte)
	}

	copy(a.Sec[:], data[:32])
	r := NewRCDe()

	pubBytes := a.PubBytes()
	if len(pubBytes) != 64 {
		return nil, fmt.Errorf("incorrect number of bytes for public key")
	}
	copy(r.Pub[:], pubBytes)
	a.RCD = r

	return data[32:], nil
}

func (a *EthSecret) MarshalBinary() ([]byte, error) {
	return a.SecBytes()[:32], nil
}

// RCDHash returns the Hash of the Redeem Condition Datastructure from a Factoid
// Address.
func (a *EthSecret) RCDHash() []byte {
	return a.RCD.Hash()
}

// RCDType returns the Redeem Condition Datastructure type used by the Factoid
// Address.
func (a *EthSecret) RCDType() uint8 {
	return a.RCD.Type()
}

// PubBytes returns the byte representation of the public key
func (a EthSecret) PubBytes() []byte {
	pub := a.PublicKey()
	bytes := crypto.FromECDSAPub(&pub)
	// Strip off the 0x04 prefix to indicate an uncompressed key.
	// You can find the prefix list here:
	// https://www.oreilly.com/library/view/mastering-ethereum/9781491971932/ch04.html
	return bytes[1:]
}

// EthAddress returns the linked ether address
func (a EthSecret) EthAddress() string {
	return crypto.PubkeyToAddress(a.PublicKey()).String()
}

// SecBytes returns the []byte representation of the secret key.
func (a *EthSecret) SecBytes() []byte {
	return a.Sec[:]
}

// SecFixed returns the fixed size secret key ([64]byte).
func (a *EthSecret) SecFixed() *[32]byte {
	return a.Sec
}

// SecString returns the string encoding of the secret key i.e. 0x123...
func (a EthSecret) SecString() string {
	str := hex.EncodeToString(a.Sec[:])
	return "0x" + str
}

// Returns the FA address
func (a EthSecret) FAString() string {
	return a.String()[1:]
}

// Returns the eFA address
func (a EthSecret) String() string {
	buf := new(bytes.Buffer)

	// FC address prefix
	buf.Write(fcPubPrefix)

	// RCD Hash
	buf.Write(a.RCDHash())

	// Checksum
	check := shad(buf.Bytes())[:ChecksumLength]
	buf.Write(check)

	return "e" + base58.Encode(buf.Bytes())
}

// ParseMnemonic parse and validate a bip39 mnumonic string. Remove extra
// spaces, capitalization, etc. Return an error if the string is invalid.
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
