// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"fmt"
	
	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
)

var (
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
	fcPubPrefix = []byte{0x5f, 0xb1}
	fcSecPrefix = []byte{0x64, 0x78}
)

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
	pub *[ed.PublicKeySize]byte
	sec *[ed.PrivateKeySize]byte
}

func NewECAddress() *ECAddress {
	a := new(ECAddress)
	a.pub = new([ed.PublicKeySize]byte)
	a.sec = new([ed.PrivateKeySize]byte)
	return a
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
	
	a := NewECAddress()
	copy(a.sec[:], p[2:34])
	// GetPublicKey will overwrite the pubkey portion of 'a.sec'
	a.pub = ed.GetPublicKey(a.sec)
	
	return a, nil
}

func (a *ECAddress) PubBytes() []byte {
	return a.pub[:]
}

func (a *ECAddress) PubFixed() *[32]byte {
	return a.pub
}

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

func (a *ECAddress) SecBytes() []byte {
	return a.sec[:]
}

func (a *ECAddress) SecFixed() *[64]byte {
	return a.sec
}

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

func (a *ECAddress) Sign(msg []byte) *[ed.SignatureSize]byte {
	return ed.Sign(a.SecFixed(), msg)
}

func (a *ECAddress) String() string {
	return a.PubString()
}

type FactoidAddress struct {
	rcd RCD
	sec *[ed.PrivateKeySize]byte
}

func NewFactoidAddress() *FactoidAddress {
	a := new(FactoidAddress)
	r := NewRCD1()
	r.Pub = new([ed.PublicKeySize]byte)
	a.rcd = r
	a.sec = new([ed.PrivateKeySize]byte)
	return a
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
	
	a := NewFactoidAddress()
	copy(a.sec[:], p[2:34])
	// GetPublicKey will overwrite the pubkey portion of 'a.sec'
	r := NewRCD1()
	r.Pub = ed.GetPublicKey(a.sec)
	a.rcd = r
	
	return a, nil
}

func (a *FactoidAddress) RCDHash() []byte {
	return a.rcd.Hash()
}

func (a *FactoidAddress) RCDType() uint8 {
	return a.rcd.Type()
}

func (a *FactoidAddress) PubString() string {
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

func (a *FactoidAddress) SecBytes() []byte {
	return a.sec[:]
}

func (a *FactoidAddress) SecFixed() *[64]byte {
	return a.sec
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
	return a.PubString()
}
