// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factom

import (
	"bytes"
	"fmt"
	"io"
	
	"github.com/FactomProject/btcutil/base58"
	ed "github.com/FactomProject/ed25519"
)

var (
	ecPubPrefix = []byte{0x59, 0x2a}
	ecSecPrefix = []byte{0x5d, 0xb6}
)

func IsValidAddress(s string) bool {
	buf := bytes.NewBuffer(base58.Decode(s))

	p := make([]byte, 34)
	if _, err := io.ReadFull(buf, p); err != nil {
		return false
	}
	
	prefix := p[:2]
	if !bytes.Equal(prefix, ecPubPrefix) && !bytes.Equal(prefix, ecSecPrefix) {
		return false
	}
	
	check := make([]byte, 4)
	if _, err := io.ReadFull(buf, check); err != nil {
		return false
	}
	
	// return true iff the checksum matches
	if bytes.Equal(shad(p)[:4], check) {
		return true
	}
	
	return false
}

type ECAddress struct {
	pub *[ed.PublicKeySize]byte
	sec *[ed.PrivateKeySize]byte
}

func NewECAddress() *ECAddress {
	e := new(ECAddress)
	e.pub = new([ed.PublicKeySize]byte)
	e.sec = new([ed.PrivateKeySize]byte)
	return e
}

// GetECAddress takes a private address 'Es...' and returns 
func GetECAddress(s string) (*ECAddress, error) {
	if !IsValidAddress(s) {
		return nil, fmt.Errorf("Invalid Address")
	}
	
	e := NewECAddress()
	copy(e.sec[:], base58.Decode(s)[2:34])
	// GetPublicKey will overwrite the pubkey portion of 'key'
	e.pub = ed.GetPublicKey(e.sec)
	
	return e, nil
}

func (e *ECAddress) IsValid() bool {
	if !IsValidAddress(e.PubString()) {
		return false
	} else if !bytes.Equal(e.pub[:2], ecPubPrefix) {
		return false
	} else if !IsValidAddress(e.SecString()) {
		return false
	} else if !bytes.Equal(e.sec[:2], ecSecPrefix) {
		return false
	} else {
		return true
	}
	// should never reach here
	return false
}

func (e *ECAddress) PubBytes() []byte {
	return e.pub[:]
}

func (e *ECAddress) SecBytes() []byte {
	return e.sec[:]
}

func (e *ECAddress) PubFixed() *[32]byte {
	return e.pub
}

func (e *ECAddress) SecFixed() *[64]byte {
	return e.sec
}

func (e *ECAddress) PubString() string {
	buf := new(bytes.Buffer)
	
	// EC address prefix
	buf.Write(ecPubPrefix)
	
	// Public key
	buf.Write(e.PubBytes())
	
	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)
	
	return base58.Encode(buf.Bytes())
}

func (e *ECAddress) SecString() string {
	buf := new(bytes.Buffer)
	
	// EC address prefix
	buf.Write(ecSecPrefix)
	
	// Secret key
	buf.Write(e.SecBytes()[:32])
	
	// Checksum
	check := shad(buf.Bytes())[:4]
	buf.Write(check)
	
	return base58.Encode(buf.Bytes())
}

func (e *ECAddress) String() string {
	return e.PubString()
}
